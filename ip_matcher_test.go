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
