package hostmetrics

import (
	"context"
	"encoding/json"
	"testing"
)

// TestCollectNetwork_IntegrationReal performs an integration test with the actual system interfaces
func TestCollectNetwork_IntegrationReal(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Test 1: No filter (default behavior)
	t.Run("NoFilter", func(t *testing.T) {
		snapshot, err := Collect(ctx, nil, "", "")
		if err != nil {
			t.Fatalf("Collect failed: %v", err)
		}

		t.Logf("Found %d network interfaces (no filter)", len(snapshot.Network))
		for _, iface := range snapshot.Network {
			t.Logf("  - %s: %v", iface.Name, iface.Addresses)
		}

		if len(snapshot.Network) == 0 {
			t.Error("Expected at least one network interface")
		}
	})

	// Test 2: Filter by first available interface
	t.Run("FilterByInterface", func(t *testing.T) {
		// First get all interfaces
		snapshot, err := Collect(ctx, nil, "", "")
		if err != nil {
			t.Fatalf("Collect failed: %v", err)
		}

		if len(snapshot.Network) == 0 {
			t.Skip("No network interfaces available for filtering test")
		}

		// Filter by the first interface name
		firstInterface := snapshot.Network[0].Name
		filteredSnapshot, err := Collect(ctx, nil, firstInterface, "")
		if err != nil {
			t.Fatalf("Collect with filter failed: %v", err)
		}

		t.Logf("Filtered to interface: %s", firstInterface)
		t.Logf("Found %d network interfaces after filtering", len(filteredSnapshot.Network))

		if len(filteredSnapshot.Network) != 1 {
			t.Errorf("Expected 1 interface after filtering by %s, got %d", firstInterface, len(filteredSnapshot.Network))
		}

		if len(filteredSnapshot.Network) > 0 && filteredSnapshot.Network[0].Name != firstInterface {
			t.Errorf("Expected interface %s, got %s", firstInterface, filteredSnapshot.Network[0].Name)
		}
	})

	// Test 3: Filter by first available IP
	t.Run("FilterByIP", func(t *testing.T) {
		// First get all interfaces
		snapshot, err := Collect(ctx, nil, "", "")
		if err != nil {
			t.Fatalf("Collect failed: %v", err)
		}

		if len(snapshot.Network) == 0 || len(snapshot.Network[0].Addresses) == 0 {
			t.Skip("No network interfaces with addresses available for filtering test")
		}

		// Get the first IP from the first interface and strip CIDR notation
		firstIP := stripCIDRSuffix(snapshot.Network[0].Addresses[0])

		filteredSnapshot, err := Collect(ctx, nil, "", firstIP)
		if err != nil {
			t.Fatalf("Collect with IP filter failed: %v", err)
		}

		t.Logf("Filtered to IP: %s", firstIP)
		t.Logf("Found %d network interfaces after filtering", len(filteredSnapshot.Network))

		if len(filteredSnapshot.Network) == 0 {
			t.Errorf("Expected at least 1 interface after filtering by IP %s, got 0", firstIP)
		}

		// Verify the filtered interface has the requested IP
		found := false
		for _, iface := range filteredSnapshot.Network {
			for _, addr := range iface.Addresses {
				if stripCIDRSuffix(addr) == firstIP {
					found = true
					break
				}
			}
		}

		if !found {
			t.Errorf("Filtered interfaces don't contain the requested IP %s", firstIP)
		}
	})

	// Test 4: Filter by non-existent interface (should fallback)
	t.Run("FilterNonExistentInterfaceFallback", func(t *testing.T) {
		snapshot, err := Collect(ctx, nil, "eth99999", "")
		if err != nil {
			t.Fatalf("Collect failed: %v", err)
		}

		t.Logf("Filtered to non-existent interface eth99999")
		t.Logf("Found %d network interfaces (should fallback to all)", len(snapshot.Network))

		// Should fall back to showing all interfaces
		if len(snapshot.Network) == 0 {
			t.Error("Expected fallback to show all interfaces, got 0")
		}
	})

	// Test 5: Output sample for documentation
	t.Run("SampleOutput", func(t *testing.T) {
		snapshot, err := Collect(ctx, nil, "", "")
		if err != nil {
			t.Fatalf("Collect failed: %v", err)
		}

		if len(snapshot.Network) > 0 {
			data, _ := json.MarshalIndent(snapshot.Network[0], "", "  ")
			t.Logf("Sample network interface output:\n%s", string(data))
		}
	})
}
