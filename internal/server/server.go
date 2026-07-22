package server

import (
	"fmt"
	"log"
	"time"

	"github.com/Austinhamilton1/dnsproxy/internal/blocker"
	"github.com/Austinhamilton1/dnsproxy/internal/cache"
	"github.com/Austinhamilton1/dnsproxy/internal/upstream"
	"github.com/miekg/dns"
)

type Server struct {
	addr    string
	blocker *blocker.Blocker
	cache   *cache.Cache
}

func New(addr string, blockedFile string) *Server {
	cache := cache.New()

	// Cleanup expired entries every minute
	go cache.Cleanup(time.Minute)

	if len(blockedFile) > 0 {
		blocker, err := blocker.Load(blockedFile)
		if err != nil {
			log.Fatalf("could not create blocker module: %s", err)
		}

		return &Server{
			addr:    addr,
			blocker: blocker,
			cache:   cache,
		}
	}

	return &Server{
		addr:    addr,
		blocker: blocker.New([]string{}),
		cache:   cache,
	}
}

func (s *Server) Run() error {
	dns.HandleFunc(".", s.handle)

	server := &dns.Server{
		Addr: s.addr,
		Net:  "udp",
	}

	fmt.Printf("Listening on %s\n", s.addr)
	return server.ListenAndServe()
}

func (s *Server) handle(
	w dns.ResponseWriter,
	r *dns.Msg,
) {
	for _, q := range r.Question {
		fmt.Println(q.Name)

		if s.blocker.IsBlocked(q.Name) {
			fmt.Printf(
				"[BLOCK] %s (%s)\n",
				q.Name,
				w.RemoteAddr(),
			)

			// Create NXDOMAIN response
			msg := new(dns.Msg)
			msg.SetReply(r)
			msg.Rcode = dns.RcodeNameError
			w.WriteMsg(msg)

			return
		}

		if msg, ok := s.cache.Get(r); ok {
			fmt.Printf(
				"[CACHE HIT] %s\n",
				q.Name,
			)

			if err := w.WriteMsg(msg); err != nil {
				fmt.Println(err)
			}

			return
		}

		fmt.Printf(
			"[CACHE MISS] %s\n",
			q.Name,
		)
	}

	response, err := upstream.Forward(r)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(r.Question) > 0 {
		s.cache.Set(r.Question[0], response)
	}

	if err := w.WriteMsg(response); err != nil {
		fmt.Println(err)
	}
}
