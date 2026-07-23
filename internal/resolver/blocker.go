package resolver

import (
	"github.com/Austinhamilton1/dnsproxy/internal/blocker"
	"github.com/miekg/dns"
)

type Blocker struct {
	next    Resolver
	blocker *blocker.Blocker
}

func NewBlocker(b *blocker.Blocker, next Resolver) *Blocker {
	return &Blocker{
		next:    next,
		blocker: b,
	}
}

func (b *Blocker) Resolve(req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) == 0 {
		return b.next.Resolve(req)
	}

	q := req.Question[0]

	if !b.blocker.IsBlocked(q.Name) {
		return b.next.Resolve(req)
	}

	msg := new(dns.Msg)
	msg.SetReply(req)
	msg.Rcode = dns.RcodeNameError

	return msg, nil
}
