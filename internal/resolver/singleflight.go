package resolver

import (
	"sync"
	"sync/atomic"

	"github.com/Austinhamilton1/dnsproxy/internal/cache"
	"github.com/Austinhamilton1/dnsproxy/internal/logger"
	"github.com/miekg/dns"
)

type call struct {
	done    chan struct{}
	waiters atomic.Uint32
	msg     *dns.Msg
	err     error
}

type SingleFlight struct {
	mu       sync.Mutex
	inflight map[string]*call
	next     Resolver
}

func NewSingleFlight(next Resolver) *SingleFlight {
	return &SingleFlight{
		next:     next,
		inflight: make(map[string]*call),
	}
}

func (s *SingleFlight) Resolve(req *dns.Msg) (*dns.Msg, error) {
	q := req.Question[0]

	key := cache.Key(q)

	s.mu.Lock()

	// Check if there's already a DNS query inflight
	if c, ok := s.inflight[key]; ok {
		s.mu.Unlock()

		c.waiters.Add(1)

		<-c.done

		if c.msg == nil {
			return nil, c.err
		}

		resp := c.msg.Copy()
		resp.Id = req.Id
		resp.Question = req.Question

		return resp, c.err
	}

	// Send one call upstream
	c := &call{
		done: make(chan struct{}),
	}
	s.inflight[key] = c
	s.mu.Unlock()

	msg, err := s.next.Resolve(req)

	c.msg = msg
	c.err = err

	close(c.done)

	// Upstream returned, finish the single flight
	s.mu.Lock()
	delete(s.inflight, key)
	s.mu.Unlock()

	if n := c.waiters.Load(); n > 0 {
		logger.Info("[SINGLEFLIGHT]", cache.Key(q), "released", n, "waiters")
	}

	if msg != nil {
		resp := msg.Copy()
		resp.Id = req.Id
		resp.Question = req.Question

		return resp, err
	}

	return nil, err
}
