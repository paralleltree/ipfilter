package ipfilter

import (
	"fmt"
	"net"
)

type IPMatcher struct {
	ranges []net.IPNet
}

// NewIPMatcher creates a new IPMatcher.
// The ranges parameter is a list of IP ranges in CIDR notation.
func NewIPMatcher(ranges []string) (*IPMatcher, error) {
	ipnets := make([]net.IPNet, 0, len(ranges))
	for _, r := range ranges {
		_, ipnet, err := net.ParseCIDR(r)
		if err != nil {
			return nil, fmt.Errorf("ParseCIDR(%q): %v", r, err)
		}
		ipnets = append(ipnets, *ipnet)
	}
	return &IPMatcher{
		ranges: ipnets,
	}, nil
}

// Match reports whether the given IP address is in the ranges.
func (m *IPMatcher) Match(addr net.IP) bool {
	for _, r := range m.ranges {
		if r.Contains(addr) {
			return true
		}
	}
	return false
}
