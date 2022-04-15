package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultWriteTimeout  = "20s"
	defaultReadTimeout   = "20s"
)

var validate *validator.Validate = validator.New()

type (
	Config struct {
		Server `toml:"Server"`
	}

	Server struct {
		ListenAddr   string `toml:"listen_addr"`
		WriteTimeout string `toml:"write_timeout"`
		ReadTimeout  string `toml:"read_timeout"`
	}
)

// Validate returns an error if the Config object is invalid.
func (c Config) Validate() error {
	return validate.Struct(c)
}

// ParseConfig attempts to read and parse configuration from the given file path.
// An error is returned if reading or parsing the config fails.
func ParseConfig(configPath string) (Config, error) {
	var cfg Config

	if configPath == "" {
		return cfg, fmt.Errorf("Empty config path")
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}

	if _, err := toml.Decode(string(configData), &cfg); err != nil {
		return cfg, fmt.Errorf("failed to decode config: %w", err)
	}

	if len(cfg.Server.WriteTimeout) == 0 {
		cfg.Server.WriteTimeout = defaultWriteTimeout
	}
	if len(cfg.Server.ReadTimeout) == 0 {
		cfg.Server.ReadTimeout = defaultReadTimeout
	}
	if len(cfg.Server.ListenAddr) == 0 {
		cfg.Server.ListenAddr = defaultServerAddress
	}

	return cfg, cfg.Validate()
}
