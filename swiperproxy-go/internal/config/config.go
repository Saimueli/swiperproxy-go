package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Proxy     ProxyConfig     `yaml:"proxy"`
	Cache     CacheConfig     `yaml:"cache"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Security  SecurityConfig  `yaml:"security"`
	Logging   LoggingConfig   `yaml:"logging"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type ProxyConfig struct {
	Target       string        `yaml:"target"`
	Timeout      time.Duration `yaml:"timeout"`
	MaxRedirects int           `yaml:"max_redirects"`
}

type CacheConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Type      string `yaml:"type"`
	MaxSize   int    `yaml:"max_size"`
	TTL       int    `yaml:"ttl"`
	RedisURL  string `yaml:"redis_url"`
}

type RateLimitConfig struct {
	Enabled            bool `yaml:"enabled"`
	RequestsPerSecond  int  `yaml:"requests_per_second"`
	Burst              int  `yaml:"burst"`
	CleanupInterval    int  `yaml:"cleanup_interval"`
}

type SecurityConfig struct {
	EnableHeaders   bool     `yaml:"enable_headers"`
	EnableCORS      bool     `yaml:"enable_cors"`
	AllowedOrigins  []string `yaml:"allowed_origins"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	File     string `yaml:"file"`
	MaxSize  int    `yaml:"max_size"`
	MaxBackups int  `yaml:"max_backups"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}