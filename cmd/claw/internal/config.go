package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/tigorlazuardi/claw/lib/claw/config"
)

var cfg *CLIConfig

type CLIConfig struct {
	Server   ServerConfig `koanf:"server"`
	Database struct {
		Path string `koanf:"path"`
	} `koanf:"database"`
	Claw *config.Config `koanf:"claw"`

	Koanf        *koanf.Koanf `koanf:"-"`
	FileLocation string       `koanf:"-"`
	Parser       koanf.Parser `koanf:"-"`
	FileProvider *file.File   `koanf:"-"`
	OnChange     func()       `koanf:"-"`
}

type ServerConfig struct {
	Host    *NetListener `koanf:"host"`
	BaseURL string       `koanf:"base_url"`
	// WebUI path to load static files from..
	// If empty, the web UI will use embedded files.
	WebUI WebUI `koanf:"webui"`
}

type WebUI struct {
	Title string `koanf:"title"`
	// Path to load static files from.
	// If empty, the web UI will use embedded files.
	Path string `koanf:"path"`
	// Enable development mode.
	// In development mode, the web UI will proxy requests to the Vite dev server.
	//
	// Unless you are developing this application, you should not enable this.
	DevMode bool `koanf:"dev_mode"`
}

func defaultCLIConfig() *CLIConfig {
	c := &CLIConfig{}
	c.Database.Path = "claw.db"
	c.Claw = config.DefaultConfig()
	c.Koanf = koanf.New(".")
	c.FileLocation = filepath.Join(xdg.ConfigHome, "claw", "config.yaml")
	c.Parser = yaml.Parser()
	return c
}

func (c *CLIConfig) LoadEnv() error {
	c.Koanf.Load(env.Provider("CLAW_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "CLAW_")
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "__", ".")
	}), nil)
	return c.Koanf.Unmarshal("", c)
}

// ReadAndWatch() sets up a file watcher on the configuration file to reload on changes.
// If OnChange is not nil, it will be called after the configuration is reloaded.
func (cfg *CLIConfig) ReadAndWatch() error {
	cfg.FileProvider = file.Provider(cfg.FileLocation)
	stat, err := os.Stat(cfg.FileLocation)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if stat != nil && stat.IsDir() {
		if stat.IsDir() {
			return fmt.Errorf("config path %q is a directory", cfg.FileLocation)
		}
		if err := cfg.Koanf.Load(cfg.FileProvider, cfg.Parser); err != nil {
			return err
		}
		if err := cfg.Koanf.Unmarshal("", cfg); err != nil {
			return err
		}
	}
	cfg.FileProvider.Watch(func(_ any, err error) {
		if err != nil {
			slog.Error("error watching config file", "error", err)
			return
		}
		if err := cfg.Koanf.Load(cfg.FileProvider, cfg.Parser); err != nil {
			slog.Error("error reloading config file", "error", err)
			return
		}
		if err := cfg.Koanf.Unmarshal("", cfg); err != nil {
			slog.Error("error unmarshaling config file", "error", err)
			return
		}
		slog.Info("configuration file reloaded", "file", cfg.FileLocation)
		if cfg.OnChange != nil {
			cfg.OnChange()
		}
	})
	return nil
}

func (cfg *CLIConfig) Unwatch() error {
	if cfg.FileProvider != nil {
		return cfg.FileProvider.Unwatch()
	}
	return nil
}
