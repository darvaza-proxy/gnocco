package gnocco

import (
	"context"
	"fmt"
	"net"
	"time"

	"darvaza.org/resolver"
	"github.com/miekg/dns"
)

// HandleRequest is the default request handler
func HandleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.RecursionAvailable = true

	q := r.Question[0]
	switch q.Qclass {
	case dns.ClassCHAOS:
		handleCHAOS(w, r)
	case dns.ClassINET:
		z, err := resolver.NewRootLookuper("")
		if err != nil {
			// Cannot build RootLookuper
			fmt.Println(err)
			m.SetRcode(r, dns.RcodeServerFailure)
			w.WriteMsg(m)
			break
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		rsp, err := z.Lookup(ctx, q.Name, q.Qtype)
		if err != nil {
			handleLookupErr(w, r, err)
			break
		}
		if rsp == nil {
			fmt.Println("Nil Answer from resolver")
			m.SetRcode(r, dns.RcodeServerFailure)
			w.WriteMsg(m)
		} else {
			rsp.SetReply(r)
			rsp.SetRcode(r, dns.RcodeSuccess)
			w.WriteMsg(rsp)
		}
	default:
		m.SetRcode(r, dns.RcodeNotImplemented)
		w.WriteMsg(m)
	}
}

func handleCHAOS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.RecursionAvailable = true

	q := r.Question[0]

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
}

func handleLookupErr(w dns.ResponseWriter, r *dns.Msg, err error) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.RecursionAvailable = true
	if n, ok := err.(*net.DNSError); ok {
		if n.Err == "NXDOMAIN" {
			m.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(m)
			return
		}
	}
	// NOTYPE and possibly others arrive here
	m.SetRcode(r, dns.RcodeSuccess)
	w.WriteMsg(m)
}
