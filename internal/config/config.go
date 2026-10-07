package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySQL    MySQLConfig    `yaml:"mysql"`
	DeepSeek DeepSeekConfig `yaml:"deepseek"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

type DeepSeekConfig struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.MySQL.DSN == "" {
		return nil, fmt.Errorf("mysql.dsn is required")
	}

	return &cfg, nil
}
