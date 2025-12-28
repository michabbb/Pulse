package hostmetrics

import "strings"

// stripCIDRSuffix removes the CIDR notation from an IP address string (e.g., "192.168.1.1/24" -> "192.168.1.1")
func stripCIDRSuffix(addr string) string {
	before, _, _ := strings.Cut(addr, "/")
	return before
}
