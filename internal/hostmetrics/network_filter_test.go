package hostmetrics

import (
	"context"
	"testing"

	gonet "github.com/shirou/gopsutil/v4/net"
)

func TestCollectNetwork_NoFilter(t *testing.T) {
	originalNetInterfaces := netInterfaces
	originalNetIOCounters := netIOCounters
	t.Cleanup(func() {
		netInterfaces = originalNetInterfaces
		netIOCounters = originalNetIOCounters
	})

	// Mock interfaces
	netInterfaces = func(ctx context.Context) (gonet.InterfaceStatList, error) {
		return gonet.InterfaceStatList{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "192.168.1.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
			{
				Name:         "eth1",
				HardwareAddr: "00:11:22:33:44:66",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "10.0.0.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
		}, nil
	}

	netIOCounters = func(ctx context.Context, pernic bool) ([]gonet.IOCountersStat, error) {
		return []gonet.IOCountersStat{
			{Name: "eth0", BytesRecv: 1000, BytesSent: 2000},
			{Name: "eth1", BytesRecv: 3000, BytesSent: 4000},
		}, nil
	}

	ctx := context.Background()
	result := collectNetwork(ctx, "", "")

	if len(result) != 2 {
		t.Errorf("expected 2 interfaces, got %d", len(result))
	}
}

func TestCollectNetwork_FilterByInterface(t *testing.T) {
	originalNetInterfaces := netInterfaces
	originalNetIOCounters := netIOCounters
	t.Cleanup(func() {
		netInterfaces = originalNetInterfaces
		netIOCounters = originalNetIOCounters
	})

	// Mock interfaces
	netInterfaces = func(ctx context.Context) (gonet.InterfaceStatList, error) {
		return gonet.InterfaceStatList{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "192.168.1.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
			{
				Name:         "eth1",
				HardwareAddr: "00:11:22:33:44:66",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "10.0.0.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
		}, nil
	}

	netIOCounters = func(ctx context.Context, pernic bool) ([]gonet.IOCountersStat, error) {
		return []gonet.IOCountersStat{
			{Name: "eth0", BytesRecv: 1000, BytesSent: 2000},
			{Name: "eth1", BytesRecv: 3000, BytesSent: 4000},
		}, nil
	}

	ctx := context.Background()
	result := collectNetwork(ctx, "eth1", "")

	if len(result) != 1 {
		t.Errorf("expected 1 interface, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "eth1" {
		t.Errorf("expected interface eth1, got %s", result[0].Name)
	}
}

func TestCollectNetwork_FilterByIP(t *testing.T) {
	originalNetInterfaces := netInterfaces
	originalNetIOCounters := netIOCounters
	t.Cleanup(func() {
		netInterfaces = originalNetInterfaces
		netIOCounters = originalNetIOCounters
	})

	// Mock interfaces
	netInterfaces = func(ctx context.Context) (gonet.InterfaceStatList, error) {
		return gonet.InterfaceStatList{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "192.168.1.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
			{
				Name:         "eth1",
				HardwareAddr: "00:11:22:33:44:66",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "10.0.0.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
		}, nil
	}

	netIOCounters = func(ctx context.Context, pernic bool) ([]gonet.IOCountersStat, error) {
		return []gonet.IOCountersStat{
			{Name: "eth0", BytesRecv: 1000, BytesSent: 2000},
			{Name: "eth1", BytesRecv: 3000, BytesSent: 4000},
		}, nil
	}

	ctx := context.Background()
	result := collectNetwork(ctx, "", "10.0.0.10")

	if len(result) != 1 {
		t.Errorf("expected 1 interface, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "eth1" {
		t.Errorf("expected interface eth1, got %s", result[0].Name)
	}
}

func TestCollectNetwork_FilterByBoth(t *testing.T) {
	originalNetInterfaces := netInterfaces
	originalNetIOCounters := netIOCounters
	t.Cleanup(func() {
		netInterfaces = originalNetInterfaces
		netIOCounters = originalNetIOCounters
	})

	// Mock interfaces
	netInterfaces = func(ctx context.Context) (gonet.InterfaceStatList, error) {
		return gonet.InterfaceStatList{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "192.168.1.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
			{
				Name:         "eth1",
				HardwareAddr: "00:11:22:33:44:66",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "10.0.0.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
		}, nil
	}

	netIOCounters = func(ctx context.Context, pernic bool) ([]gonet.IOCountersStat, error) {
		return []gonet.IOCountersStat{
			{Name: "eth0", BytesRecv: 1000, BytesSent: 2000},
			{Name: "eth1", BytesRecv: 3000, BytesSent: 4000},
		}, nil
	}

	ctx := context.Background()
	result := collectNetwork(ctx, "eth1", "10.0.0.10")

	if len(result) != 1 {
		t.Errorf("expected 1 interface, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "eth1" {
		t.Errorf("expected interface eth1, got %s", result[0].Name)
	}
}

func TestCollectNetwork_FilterNotFoundFallback(t *testing.T) {
	originalNetInterfaces := netInterfaces
	originalNetIOCounters := netIOCounters
	t.Cleanup(func() {
		netInterfaces = originalNetInterfaces
		netIOCounters = originalNetIOCounters
	})

	// Mock interfaces
	netInterfaces = func(ctx context.Context) (gonet.InterfaceStatList, error) {
		return gonet.InterfaceStatList{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "192.168.1.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
		}, nil
	}

	netIOCounters = func(ctx context.Context, pernic bool) ([]gonet.IOCountersStat, error) {
		return []gonet.IOCountersStat{
			{Name: "eth0", BytesRecv: 1000, BytesSent: 2000},
		}, nil
	}

	ctx := context.Background()
	// Try to filter by non-existent interface - should fall back to all interfaces
	result := collectNetwork(ctx, "eth99", "")

	if len(result) != 1 {
		t.Errorf("expected fallback to 1 interface, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "eth0" {
		t.Errorf("expected fallback to eth0, got %s", result[0].Name)
	}
}

func TestCollectNetwork_LoopbackFiltered(t *testing.T) {
	originalNetInterfaces := netInterfaces
	originalNetIOCounters := netIOCounters
	t.Cleanup(func() {
		netInterfaces = originalNetInterfaces
		netIOCounters = originalNetIOCounters
	})

	// Mock interfaces including loopback
	netInterfaces = func(ctx context.Context) (gonet.InterfaceStatList, error) {
		return gonet.InterfaceStatList{
			{
				Name:         "lo",
				HardwareAddr: "",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "127.0.0.1/8"},
				},
				Flags: []string{"up", "loopback"},
			},
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				Addrs: []gonet.InterfaceAddr{
					{Addr: "192.168.1.10/24"},
				},
				Flags: []string{"up", "broadcast", "multicast"},
			},
		}, nil
	}

	netIOCounters = func(ctx context.Context, pernic bool) ([]gonet.IOCountersStat, error) {
		return []gonet.IOCountersStat{
			{Name: "lo", BytesRecv: 1000, BytesSent: 2000},
			{Name: "eth0", BytesRecv: 3000, BytesSent: 4000},
		}, nil
	}

	ctx := context.Background()
	result := collectNetwork(ctx, "", "")

	// Loopback should be filtered out automatically
	if len(result) != 1 {
		t.Errorf("expected 1 non-loopback interface, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "eth0" {
		t.Errorf("expected interface eth0, got %s", result[0].Name)
	}
}
