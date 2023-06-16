package gnocco

import (
	"context"

	"darvaza.org/resolver"
	"github.com/miekg/dns"
)

// HandleRequest is the default request handler
// revive:disable:cyclomatic
// revive:disable:cognitive-complexity
func HandleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	q := r.Question[0]
	switch {
	case q.Qclass == dns.ClassCHAOS:
		hdr := dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
		switch q.Name {
		case "authors.bind.":
			m.Answer = append(m.Answer, &dns.TXT{Hdr: hdr,
				Txt: []string{"JPI Technologies <oss@jpi.io>"}})
			m.SetRcode(r, dns.RcodeSuccess)
			w.WriteMsg(m)
		case "version.bind.", "version.server.":
			m.Answer = []dns.RR{&dns.TXT{Hdr: hdr,
				Txt: []string{"Version " + Version + " built on " + BuildDate}}}
			m.SetRcode(r, dns.RcodeSuccess)
			w.WriteMsg(m)
		case "hostname.bind.", "id.server.":
			m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{"localhost"}}}
			m.SetRcode(r, dns.RcodeSuccess)
			w.WriteMsg(m)
		default:
			m.SetRcode(r, dns.RcodeNotImplemented)
			w.WriteMsg(m)
		}
	default:
		z, err := resolver.NewRootResolver("")
		if err != nil {
			m.SetRcode(r, dns.RcodeServerFailure)
			w.WriteMsg(m)
		}
		switch q.Qtype {
		case dns.TypeA, dns.TypeAAAA:
			resp, err := z.LookupIPAddr(context.TODO(), q.Name)
			if err != nil {
				m.SetRcode(r, dns.RcodeServerFailure)
				w.WriteMsg(m)
			}
			for _, r := range resp {
				rec, _ := dns.NewRR(dns.Fqdn(q.Name) + " " +
					dns.TypeToString[q.Qtype] + " " + r.IP.String())
				if rec != nil { // we get mixed IPv4 and IPv6 so please no nil
					m.Answer = append(m.Answer, rec)
				}
			}
			m.SetRcode(r, dns.RcodeSuccess)
			w.WriteMsg(m)
		}
	}
}
