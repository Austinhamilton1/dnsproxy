package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type DNSConfig struct {
	Listen string `toml:"listen"`
}

type CacheConfig struct {
	CleanupInterval int `toml:"cleanup_interval"`
}

type UpstreamConfig struct {
	Server string `toml:"server"`
}

type BlocklistConfig struct {
	File string `toml:"file"`
}

type LogConfig struct {
	Level int `toml:"level"`
}

type Config struct {
	DNS       DNSConfig       `toml:"dns"`
	Cache     CacheConfig     `toml:"cache"`
	Upstream  UpstreamConfig  `toml:"upstream"`
	Blocklist BlocklistConfig `toml:"blocklist"`
	Log       LogConfig       `toml:"log"`
}

func Load(filename string) (*Config, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var cfg Config

	err = toml.Unmarshal(bytes, &cfg)

	return &cfg, err
}
