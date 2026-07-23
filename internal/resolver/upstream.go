package resolver

import (
	"time"

	"github.com/miekg/dns"
)

type Upstream struct {
	client dns.Client
	addr   string
}

func NewUpstream(addr string) *Upstream {
	return &Upstream{
		addr: addr,
	}
}

func (u *Upstream) Resolve(req *dns.Msg) (*dns.Msg, error) {
	time.Sleep(2 * time.Second)
	resp, _, err := u.client.Exchange(req, u.addr)

	return resp, err
}
