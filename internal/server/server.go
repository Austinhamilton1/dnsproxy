package server

import (
	"github.com/Austinhamilton1/dnsproxy/internal/resolver"
	"github.com/miekg/dns"
)

type Server struct {
	addr     string
	resolver resolver.Resolver
}

func New(addr string, resolver resolver.Resolver) *Server {
	return &Server{
		addr:     addr,
		resolver: resolver,
	}
}

func (s *Server) Run() error {
	dns.HandleFunc(".", s.handle)

	server := &dns.Server{
		Addr: s.addr,
		Net:  "udp",
	}

	return server.ListenAndServe()
}

func (s *Server) handle(
	w dns.ResponseWriter,
	r *dns.Msg,
) {
	resp, err := s.resolver.Resolve(r)

	if err != nil {
		return
	}

	_ = w.WriteMsg(resp)
}
