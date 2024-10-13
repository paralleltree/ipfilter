package ipfilter_test

import (
	"net"
	"testing"

	"github.com/paralleltree/ipfilter"
)

func TestIPMatcher_Match(t *testing.T) {
	var tests = []struct {
		name       string
		addr       string
		ranges     []string
		wantResult bool
	}{
		{
			name:       "single IPv4 range that matches",
			addr:       "192.168.0.1",
			ranges:     []string{"192.168.0.0/24"},
			wantResult: true,
		},
		{
			name:       "single IPv4 range that does not match",
			addr:       "192.168.1.1",
			ranges:     []string{"192.168.0.0/24"},
			wantResult: false,
		},
		{
			// addr[3]  = 0b10000001
			// range[3] = 0b10000000
			name:       "single IPv4 range that matches",
			addr:       "192.168.1.129",
			ranges:     []string{"192.168.1.128/26"},
			wantResult: true,
		},
		{
			// addr[3]  = 0b01000001
			// range[3] = 0b10000000
			name:       "single IPv4 range that does not match",
			addr:       "192.168.1.65",
			ranges:     []string{"192.168.1.128/26"},
			wantResult: false,
		},
		{
			name:       "single IPv6 range that matches",
			addr:       "2001:db8::1",
			ranges:     []string{"2001:db8::/64"},
			wantResult: true,
		},
		{
			name:       "single IPv6 range that does not match",
			addr:       "2001:db9::1",
			ranges:     []string{"2001:db8::/64"},
			wantResult: false,
		},
		{
			// the head of the address is the same as the range, but the address should not match the range
			// because they are different addresses.
			// note: the first bytes of the address and range are is [32 0 0 0].
			name:       "IPv4 range does not match IPv6 address",
			addr:       "2000::1",
			ranges:     []string{"32.0.0.0/8"},
			wantResult: false,
		},
		{
			name:       "IPv6 range does not match IPv4 address",
			addr:       "32.0.0.1",
			ranges:     []string{"2000::/8"},
			wantResult: false,
		},
		{
			name: "multiple IPv4 and IPv6 ranges that match IPv4 address",
			addr: "192.168.100.1",
			ranges: []string{
				"2001:db8::/64",
				// the ranges should be sorted internally to apply the binary search.
				"192.168.200.0/24",
				"192.168.100.0/24",
				"192.168.1.0/24",
				"192.168.0.0/24",
			},
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, err := ipfilter.NewIPMatcher(tt.ranges)
			if err != nil {
				t.Fatalf("NewIPMatcher: %v", err)
			}
			ip := net.ParseIP(tt.addr)
			gotResult := matcher.Match(ip)
			if tt.wantResult != gotResult {
				t.Errorf("unexpected result: want %v, but got %v", tt.wantResult, gotResult)
			}
		})
	}
}

func BenchmarkIPMatcher_Test(b *testing.B) {
	// on commit hash 6de25fd, the benchmark result is 1059579 ns/op
	// on commit hash 5985bca, the benchmark result is 344.5 ns/op
	n := 100000
	ranges := make([]string, 0, n+1)
	for i := 0; i < n; i++ {
		ranges = append(ranges, "192.168.0.0/24")
	}
	ranges = append(ranges, "192.168.1.0/24")
	matcher, err := ipfilter.NewIPMatcher(ranges)
	if err != nil {
		b.Fatalf("new ip matcher: %v", err)
	}
	ip := net.ParseIP("192.168.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.Match(ip)
	}
}
