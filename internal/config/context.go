package config

import (
	"os"
	"time"
)

type Context struct {
	Host         string    `yaml:"host"`
	Token        string    `yaml:"token,omitempty"`
	RefreshToken string    `yaml:"refresh_token,omitempty"`
	User         string    `yaml:"user,omitempty"`
	ExpiresAt    time.Time `yaml:"expires_at,omitempty"`
}

func (c *Context) EffectiveToken() string {
	if t := os.Getenv("ELEMENTS_TOKEN"); t != "" {
		return t
	}
	return c.Token
}

// APIURL returns the full API base URL.
// DM platform uses a single host for both API and auth, no subdomain splitting.
func (c *Context) APIURL() string {
	return "https://" + c.Host
}

func (cfg *Config) SetContext(name string, ctx *Context) {
	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]*Context)
	}
	cfg.Contexts[name] = ctx
}

func (cfg *Config) DeleteContext(name string) {
	delete(cfg.Contexts, name)
	if cfg.CurrentContext == name {
		cfg.CurrentContext = ""
	}
}
