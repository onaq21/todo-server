package config

import (
	"time"
	"os"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 				string	`yaml:"env" env-default:"local"`		
	StoragePath string	`yaml:"storage_path" env-required:"true"`
	HTTPServer					`yaml:"http_server"`
}

type HTTPServer struct {
	Address 			string			 `yaml:"address" env-default:"localhost:5000"`
	Timeout 		time.Duration	 `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration	 `yaml:"idle_timeout" env-default:"10s"`
}

func Load() (*Config, error) {
	const fn = "internal.config.Load"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &cfg, nil
}