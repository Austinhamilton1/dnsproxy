package resolver

import (
	"github.com/Austinhamilton1/dnsproxy/internal/cache"
	"github.com/Austinhamilton1/dnsproxy/internal/logger"
	"github.com/miekg/dns"
)

type Cache struct {
	next  Resolver
	cache *cache.Cache
}

func NewCache(c *cache.Cache, next Resolver) *Cache {
	return &Cache{
		next:  next,
		cache: c,
	}
}

func (c *Cache) Resolve(req *dns.Msg) (*dns.Msg, error) {
	q := req.Question[0]

	if msg, ok := c.cache.Get(req); ok {
		logger.Info("[CACHE HIT]", cache.Key(q))
		return msg, nil
	}

	logger.Info("[CACHE MISS]", cache.Key(q))

	msg, err := c.next.Resolve(req)

	if err != nil {
		return nil, err
	}

	c.cache.Set(q, msg)

	return msg, nil
}
