package ipfilter

import (
	"fmt"
	"net"
	"slices"
)

type IPMatcher struct {
	ipv4Ranges ipRangeSet
	ipv6Ranges ipRangeSet
}

// ipRangeSet is a map of IP address ranges grouped by the mask length.
type ipRangeSet map[int][]net.IP

// contains reports whether the IP address is in the range set.
// The IP address must be the same length as the range.
func (r ipRangeSet) contains(ip net.IP) bool {
	for length, ips := range r {
		_, ok := slices.BinarySearchFunc(ips, ip, func(i, j net.IP) int {
			return compareIP(length, i, j)
		})
		if ok {
			return true
		}
	}
	return false
}

// sort sorts the IP address ranges in the set.
func (r ipRangeSet) sort() {
	for length, ips := range r {
		slices.SortFunc(ips, func(i, j net.IP) int {
			return compareIP(length, i, j)
		})
	}
}

func newIPRangeSetFromRangeString(ranges []string) (ipRangeSet, ipRangeSet, error) {
	ipv4Range := ipRangeSet{}
	ipv6Range := ipRangeSet{}
	for _, r := range ranges {
		ip, ipnet, _ := net.ParseCIDR(r)
		maskLength, _ := ipnet.Mask.Size()
		if maskLength < 0 {
			return nil, nil, fmt.Errorf("invalid mask: %q", r)
		}

		if x := ip.To4(); x != nil {
			ipv4Range[maskLength] = append(ipv4Range[maskLength], x)
		} else {
			ipv6Range[maskLength] = append(ipv6Range[maskLength], ipnet.IP)
		}
	}
	// sort for binary search
	ipv4Range.sort()
	ipv6Range.sort()
	return ipv4Range, ipv6Range, nil
}

// NewIPMatcher creates a new IPMatcher.
// The ranges parameter is a list of IP ranges in CIDR notation.
func NewIPMatcher(ranges []string) (*IPMatcher, error) {
	ipv4ranges, ipv6ranges, err := newIPRangeSetFromRangeString(ranges)
	if err != nil {
		return nil, fmt.Errorf("create IPMatcher: %w", err)
	}
	return &IPMatcher{
		ipv4Ranges: ipv4ranges,
		ipv6Ranges: ipv6ranges,
	}, nil
}

// Match reports whether the given IP address is in the ranges.
func (m *IPMatcher) Match(addr net.IP) bool {
	if ipv4Addr := addr.To4(); ipv4Addr != nil {
		return m.ipv4Ranges.contains(ipv4Addr)
	}
	return m.ipv6Ranges.contains(addr)
}

// compareIP compares two IP addresses with the mask length.
func compareIP(maskLength int, a, b net.IP) int {
	if x := a.To4(); x != nil {
		a = x
	}
	if x := b.To4(); x != nil {
		b = x
	}
	if len(a) != len(b) {
		panic("ip address length not equal")
	}
	mask := net.CIDRMask(maskLength, len(a)*8)
	for i := 0; i < len(a); i++ {
		av := a[i] & mask[i]
		bv := b[i] & mask[i]
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
	}
	return 0
}
