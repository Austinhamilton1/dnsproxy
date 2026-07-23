package main

import (
	"flag"
	"time"

	"github.com/Austinhamilton1/dnsproxy/internal/blocker"
	"github.com/Austinhamilton1/dnsproxy/internal/cache"
	"github.com/Austinhamilton1/dnsproxy/internal/config"
	"github.com/Austinhamilton1/dnsproxy/internal/logger"
	"github.com/Austinhamilton1/dnsproxy/internal/resolver"
	"github.com/Austinhamilton1/dnsproxy/internal/server"
)

func main() {
	logger.Init()

	configFilePtr := flag.String("config", "", "points to the config (.toml) file for the proxy")

	flag.Parse()

	cfg, err := config.Load(*configFilePtr)
	if err != nil {
		logger.Error("could not parse config file:", err.Error())
	}

	logger.SetLevel = logger.Level(cfg.Log.Level)
	connStr := cfg.DNS.Listen

	c := cache.New()
	go c.Cleanup(time.Duration(cfg.Cache.CleanupInterval) * time.Minute)

	var r resolver.Resolver

	r = resolver.NewUpstream(cfg.Upstream.Server)

	r = resolver.NewCache(c, r)

	if cfg.Blocklist.File != "" {
		b, err := blocker.Load(cfg.Blocklist.File)
		if err != nil {
			logger.Error(err)
		}

		r = resolver.NewBlocker(b, r)
	}

	r = resolver.NewLogger(r)

	s := server.New(connStr, r)

	if err := s.Run(); err != nil {
		logger.Error("could not create DNS proxy:", err.Error())
	}
}
