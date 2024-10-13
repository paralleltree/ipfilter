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
			name:       "single range that matches",
			addr:       "192.168.0.1",
			ranges:     []string{"192.168.0.0/24"},
			wantResult: true,
		},
		{
			name:       "single range that does not match",
			addr:       "192.168.1.1",
			ranges:     []string{"192.168.0.0/24"},
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
