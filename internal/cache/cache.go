package cache

import (
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type Entry struct {
	Msg        *dns.Msg
	InsertedAt time.Time
	ExpiresAt  time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

func normalize(domain string) string {
	domain = strings.ToLower(domain)

	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	return domain
}

func minTTL(msg *dns.Msg) uint32 {
	if len(msg.Answer) == 0 {
		return 0
	}

	min := msg.Answer[0].Header().Ttl

	for _, rr := range msg.Answer[1:] {
		if rr.Header().Ttl < min {
			min = rr.Header().Ttl
		}
	}

	return min
}

func adjustTTL(rrs []dns.RR, elapsed uint32) {
	for _, rr := range rrs {
		h := rr.Header()

		if h.Ttl > elapsed {
			h.Ttl -= elapsed
		} else {
			h.Ttl = 0
		}
	}
}

func New() *Cache {
	return &Cache{
		entries: make(map[string]Entry),
	}
}

func Key(q dns.Question) string {
	return normalize(q.Name) +
		":" + dns.TypeToString[q.Qtype] +
		":" + dns.ClassToString[q.Qclass]
}

func (c *Cache) Get(req *dns.Msg) (*dns.Msg, bool) {
	// Calculate key
	key := Key(req.Question[0])

	// Lookup entry
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	// No -> Miss
	if !ok {
		return nil, false
	}

	// Expired -> Delete -> Miss
	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()

		return nil, false
	}

	msg := entry.Msg.Copy()

	msg.Id = req.Id
	msg.Question = req.Question

	elapsed := uint32(time.Since(entry.InsertedAt).Seconds())

	adjustTTL(msg.Answer, elapsed)
	adjustTTL(msg.Ns, elapsed)
	adjustTTL(msg.Extra, elapsed)

	return msg, true
}

func (c *Cache) RemoveExpired() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entries {
		if now.After(entry.ExpiresAt) {
			delete(c.entries, key)
		}
	}
}

func (c *Cache) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.RemoveExpired()
	}
}

func (c *Cache) Set(q dns.Question, msg *dns.Msg) {
	ttl := minTTL(msg)

	if ttl == 0 {
		return
	}

	now := time.Now()

	entry := Entry{
		Msg:        msg.Copy(),
		InsertedAt: now,
		ExpiresAt:  now.Add(time.Duration(ttl) * time.Second),
	}

	key := Key(q)

	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}
