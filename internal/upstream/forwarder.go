package upstream

import "github.com/miekg/dns"

func Forward(msg *dns.Msg) (*dns.Msg, error) {
	client := dns.Client{}

	response, _, err := client.Exchange(msg, "1.1.1.1:53")

	return response, err
}
