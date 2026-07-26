package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Listen string `yaml:"listen"`
	} `yaml:"server"`
	Proxy struct {
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"proxy"`
	Cache struct {
		TTL        time.Duration `yaml:"ttl"`
		MaxEntries int           `yaml:"max_entries"`
	} `yaml:"cache"`
	RateLimit struct {
		Requests int           `yaml:"requests"`
		Window   time.Duration `yaml:"window"`
	} `yaml:"rate_limit"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Cache.MaxEntries <= 0 {
		cfg.Cache.MaxEntries = 1000
	}
	return cfg, nil
}