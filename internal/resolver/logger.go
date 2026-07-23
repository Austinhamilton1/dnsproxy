package resolver

import (
	"time"

	"github.com/Austinhamilton1/dnsproxy/internal/logger"
	"github.com/miekg/dns"
)

type Logger struct {
	next Resolver
}

func NewLogger(next Resolver) *Logger {
	return &Logger{
		next: next,
	}
}

func (l *Logger) Resolve(req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) == 0 {
		return l.next.Resolve(req)
	}

	domain := req.Question[0].Name

	now := time.Now()
	res, err := l.next.Resolve(req)
	elapsed := time.Now().Sub(now).Milliseconds()

	if err != nil {
		logger.Warning("ran into error:", err.Error())
	} else {
		logger.Info("resolved:", domain, "time:", elapsed)
	}

	return res, err
}
