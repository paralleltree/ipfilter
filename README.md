# ipfilter

`ipfilter` provides a matcher to filter IP addresses based on a list of CIDR ranges.
The matcher is optimized for speed and can handle a large number of ranges efficiently.

## Installation

    $ go get github.com/paralleltree/ipfilter

## Usage

```go
import (
  "fmt"
  "net"

  "github.com/paralleltree/ipfilter"
)

func example() {
  matcher, err := ipfilter.NewIPMatcher([]string{
    "192.168.0.0/24",
    "192.168.1.0/24",
  })
  if err != nil {
    panic(err)
  }

  ips := []string{
    "192.168.0.10",
    "192.168.20.10",
  }
  for _, ip := range ips {
    result := matcher.Match(net.ParseIP(ip))
    fmt.Printf("%s: %v\n", ip, result)
  }
  // Output:
  // 192.168.0.10: true
  // 192.168.20.10: false
}
```

## How it works
The matcher groups the given CIDR ranges by their prefix(subnet mask) length.
When matching an IP address, it tries to find the corresponding range from each group using a binary search.

```mermaid
flowchart LR
  root(IPMatcher)
  ipv4set[IPv4RangeSet]
  ipv6set[IPv6RangeSet]
  ipv4PrefixLength16[prefixLength=16]
  ipv4PrefixLength24[prefixLength=24]
  ipv6PrefixLength64[prefixLength=64]
  ipv4PrefixLength24Range1([192.168.0.0/24])
  ipv4PrefixLength24Range2([192.168.1.0/24])
  ipv6PrefixLength16Range1([172.16.0.0/16])
  ipv6PrefixLength64Range1([2001:db8::/64])
  ipv6PrefixLength64Range2([78ab:920::/64])

  root --- ipv4set
  root --- ipv6set
  ipv4set --- ipv4PrefixLength16
  ipv4set --- ipv4PrefixLength24
  ipv4PrefixLength16 --- ipv6PrefixLength16Range1
  ipv4PrefixLength24 --- ipv4PrefixLength24Range1
  ipv4PrefixLength24 --- ipv4PrefixLength24Range2
  ipv6set --- ipv6PrefixLength64
  ipv6PrefixLength64 --- ipv6PrefixLength64Range1
  ipv6PrefixLength64 --- ipv6PrefixLength64Range2
```
