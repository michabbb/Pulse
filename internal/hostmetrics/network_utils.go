package hostmetrics

// stripCIDRSuffix removes the CIDR notation from an IP address string (e.g., "192.168.1.1/24" -> "192.168.1.1")
func stripCIDRSuffix(addr string) string {
	for idx := 0; idx < len(addr); idx++ {
		if addr[idx] == '/' {
			return addr[:idx]
		}
	}
	return addr
}
