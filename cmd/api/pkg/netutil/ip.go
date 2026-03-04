package netutil

import (
	"net"
	"net/netip"
)

type IPChecker struct {
	allowed []netip.Prefix
}

func NewIPChecker(prefixes []netip.Prefix) *IPChecker {
	return &IPChecker{prefixes}
}

func (c *IPChecker) Contains(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}

	for _, prefix := range c.allowed {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
