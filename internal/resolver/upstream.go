package resolver

import (
	"errors"

	"github.com/miekg/dns"
)

type Upstream struct {
	client dns.Client
	addrs  []string
}

func NewUpstream(addrs []string) *Upstream {
	return &Upstream{
		addrs: addrs,
	}
}

func (u *Upstream) Resolve(req *dns.Msg) (*dns.Msg, error) {
	for _, server := range u.addrs {
		resp, _, err := u.client.Exchange(req, server)
		if err == nil && resp != nil && resp.Rcode == dns.RcodeSuccess {
			return resp, err
		}
	}

	return nil, errors.New("could not resolve")
}
