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
	Authority    string    `yaml:"authority,omitempty"`
	ExpiresAt    time.Time `yaml:"expires_at,omitempty"`
	ClientID     string    `yaml:"client_id,omitempty"`
	ClientSecret string    `yaml:"client_secret,omitempty"`

	// Impersonate: backup of admin token while impersonating another user
	AdminToken        string    `yaml:"admin_token,omitempty"`
	AdminRefreshToken string    `yaml:"admin_refresh_token,omitempty"`
	AdminExpiresAt    time.Time `yaml:"admin_expires_at,omitempty"`
	OrgID             string    `yaml:"org_id,omitempty"`
}

// IsImpersonating returns true if currently impersonating another user.
func (c *Context) IsImpersonating() bool {
	return c.AdminToken != ""
}

func (c *Context) EffectiveToken() string {
	if t := os.Getenv("DEVICEMANAGER_TOKEN"); t != "" {
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
