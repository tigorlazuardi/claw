package config

import (
	"log/slog"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Download  Download  `koanf:"download"`
	Scheduler Scheduler `koanf:"scheduler"`
	Webhooks  Webhooks  `koanf:"webhooks"`

	OnConfigChange func(newCfg *Config) `koanf:"-"`
	koanf          *koanf.Koanf         `koanf:"-"`
	fileLocation   string               `koanf:"-"`
	parser         koanf.Parser         `koanf:"-"`
	fileProvider   *file.File           `koanf:"-"`
}

// Close closes any resources held by the Config, such as file watchers.
func (cfg *Config) Close() error {
	if cfg.fileProvider != nil {
		return cfg.fileProvider.Unwatch()
	}
	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Download:  DefaultDownload(),
		Scheduler: DefaultScheduler(),
		koanf:     koanf.New("."),
	}
}

func LoadConfigFile(path string, parser koanf.Parser) (*Config, error) {
	cfg := DefaultConfig()
	k := koanf.New(".")
	provider := file.Provider(path)
	if err := k.Load(provider, parser); err != nil {
		return nil, err
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}
	cfg.koanf = k
	cfg.fileLocation = path
	cfg.parser = parser
	cfg.fileProvider = provider
	cfg.fileProvider.Watch(func(_ any, err error) {
		if err != nil {
			return
		}
		newCfg := DefaultConfig()
		if err := newCfg.koanf.Load(cfg.fileProvider, parser); err != nil {
			return
		}
		if err := cfg.koanf.Unmarshal("", newCfg); err != nil {
			return
		}
		slog.Info("configuration file reloaded", "file", path)
		newCfg.parser = cfg.parser
		newCfg.fileLocation = cfg.fileLocation
		newCfg.fileProvider = cfg.fileProvider
		newCfg.OnConfigChange = cfg.OnConfigChange
		if newCfg.OnConfigChange != nil {
			newCfg.OnConfigChange(newCfg)
		}
		*cfg = *newCfg
	})
	return cfg, nil
}
