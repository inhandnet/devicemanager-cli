package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrentContext string              `yaml:"current-context"`
	Contexts       map[string]*Context `yaml:"contexts"`
}

func (cfg *Config) ActiveContext() (*Context, error) {
	name := os.Getenv("ELEMENTS_CONTEXT")
	if name == "" {
		name = cfg.CurrentContext
	}
	if name == "" {
		return nil, fmt.Errorf("no active context; run 'elements auth login' or 'elements config use-context <name>'")
	}
	ctx, ok := cfg.Contexts[name]
	if !ok {
		return nil, fmt.Errorf("context %q not found in config", name)
	}
	if h := os.Getenv("ELEMENTS_HOST"); h != "" {
		ctx.Host = h
	}
	return ctx, nil
}

func (cfg *Config) ActiveContextName() string {
	if name := os.Getenv("ELEMENTS_CONTEXT"); name != "" {
		return name
	}
	return cfg.CurrentContext
}

func DefaultPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "elements", "config.yaml")
}

func Load(path string) (*Config, error) {
	cfg := &Config{Contexts: make(map[string]*Context)}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]*Context)
	}
	return cfg, nil
}

func Save(cfg *Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}
