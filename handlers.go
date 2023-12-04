package gnocco

import (
	"time"

	"github.com/miekg/dns"

	"darvaza.org/resolver"
	"darvaza.org/resolver/pkg/server"
)

// NewHandler allocates an iterative DNS resolver optionally
// using the given root server as starting point.
func NewHandler(start string) (dns.Handler, error) {
	z, err := resolver.NewRootLookuper(start)
	if err != nil {
		return nil, err
	}

	h := &server.Handler{
		// CHAOS
		Authors:  "JPI Technologies <oss@jpi.io>",
		Version:  "Version " + Version + " built on " + BuildDate,
		Hostname: "localhost",
		// Lookup
		Lookuper: z,
		Timeout:  5 * time.Second,
	}

	return h, nil
}
