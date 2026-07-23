package resolver

import (
	"github.com/Austinhamilton1/dnsproxy/internal/cache"
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
	if msg, ok := c.cache.Get(req); ok {
		return msg, nil
	}

	msg, err := c.next.Resolve(req)

	if err != nil {
		return nil, err
	}

	c.cache.Set(req.Question[0], msg)

	return msg, nil
}
