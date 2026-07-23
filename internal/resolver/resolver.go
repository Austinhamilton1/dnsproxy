package resolver

import "github.com/miekg/dns"

type Resolver interface {
	Resolve(*dns.Msg) (*dns.Msg, error)
}
