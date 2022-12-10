package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DSN         string `yaml:"database_source_name"`
	MaxDongSize int    `yaml:"max_dong_size"`
	Token       string `yaml:"token"`
}

var (
	cfg    *Config
	cfgErr error
	once   sync.Once
)

func GetConfig() (Config, error) {
	once.Do(func() {
		var path string
		flag.StringVar(&path, "config", "", "path to config file")
		flag.Parse()

		if path == "" {
			cfgErr = fmt.Errorf("path to config is empty")
			return
		}

		body, err := os.ReadFile(path)
		if err != nil {
			cfgErr = err
			return
		}

		if err = yaml.Unmarshal(body, &cfg); err != nil {
			cfgErr = err
			return
		}
	})

	return *cfg, cfgErr
}
