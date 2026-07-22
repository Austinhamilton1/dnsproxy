package blocker

import (
	"bufio"
	"os"
	"strings"
)

type Blocker struct {
	blocked map[string]struct{}
}

func New(domains []string) *Blocker {
	b := &Blocker{
		blocked: make(map[string]struct{}),
	}

	for _, d := range domains {
		b.blocked[normalize(d)] = struct{}{}
	}

	return b
}

func Load(filename string) (*Blocker, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var domains []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		domains = append(domains, line)
	}

	return New(domains), scanner.Err()
}

func (b *Blocker) IsBlocked(domain string) bool {
	domain = normalize(domain)

	for {
		if _, ok := b.blocked[domain]; ok {
			return true
		}

		i := strings.Index(domain, ".")
		if i == -1 || i == len(domain)-1 {
			break
		}

		domain = domain[i+1:]
	}

	return false
}

func normalize(domain string) string {
	domain = strings.ToLower(domain)

	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	return domain
}
