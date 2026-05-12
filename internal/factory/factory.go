package factory

import (
	"net/http"
	"sync"
	"time"

	"github.com/inhandnet/devicemanager-cli/internal/api"
	"github.com/inhandnet/devicemanager-cli/internal/config"
	"github.com/inhandnet/devicemanager-cli/internal/debug"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

type Factory struct {
	IO         *iostreams.IOStreams
	ConfigPath string
	Verbose    int

	configOnce sync.Once
	config     *config.Config
	configErr  error
}

func New() *Factory {
	return &Factory{
		IO:         iostreams.System(),
		ConfigPath: config.DefaultPath(),
		Verbose:    100,
	}
}

func (f *Factory) Config() (*config.Config, error) {
	f.configOnce.Do(func() {
		f.config, f.configErr = config.Load(f.ConfigPath)
	})
	return f.config, f.configErr
}

func (f *Factory) ReloadConfig() {
	f.configOnce = sync.Once{}
	f.config = nil
	f.configErr = nil
}

func (f *Factory) SaveConfig() error {
	if f.config == nil {
		return nil
	}
	return config.Save(f.config, f.ConfigPath)
}

func (f *Factory) APIClient() (*api.APIClient, error) {
	actx, err := f.activeContext()
	if err != nil {
		return nil, err
	}
	f.debugConfig(actx)
	return api.NewAPIClient(actx.APIURL(), f.newTransport(actx), f.Verbose), nil
}

func (f *Factory) activeContext() (*config.Context, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}
	return cfg.ActiveContext()
}

func (f *Factory) debugConfig(ctx *config.Context) {
	if !debug.Enabled {
		return
	}
	cfg, _ := f.Config()
	debug.Log("context: %s", cfg.ActiveContextName())
	debug.Log("api: %s", ctx.APIURL())
	if ctx.User != "" {
		debug.Log("user: %s", ctx.User)
	}
}

func (f *Factory) newTransport(ctx *config.Context) *api.TokenTransport {
	return &api.TokenTransport{
		Token:        ctx.EffectiveToken(),
		RefreshToken: ctx.RefreshToken,
		Host:         ctx.APIURL(),
		ClientID:     ctx.ClientID,
		ClientSecret: ctx.ClientSecret,
		OnRefresh: func(accessToken, refreshToken string, expiry time.Time) {
			ctx.Token = accessToken
			if refreshToken != "" {
				ctx.RefreshToken = refreshToken
			}
			if !expiry.IsZero() {
				ctx.ExpiresAt = expiry
			}
			_ = f.SaveConfig()
		},
		Base: http.DefaultTransport,
	}
}
