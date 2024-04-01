package utils

import "github.com/miekg/dns"

func NoRecordAnswer(m *dns.Msg) *dns.Msg {
	m.Rcode = dns.RcodeNameError
	return m
}
