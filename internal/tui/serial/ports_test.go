/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package serial

import (
	"testing"
)

func TestGetPortsList(t *testing.T) {
	// This test just verifies the function doesn't panic
	// The actual ports returned depend on the system
	ports, err := GetPortsList()

	// We don't require ports to exist, just that the function works
	if err != nil {
		// Some systems may return an error if no serial support
		t.Logf("GetPortsList() returned error (may be expected): %v", err)
		return
	}

	// If no error, ports should be a valid slice (even if empty)
	if ports == nil {
		t.Error("GetPortsList() returned nil slice without error")
	}

	t.Logf("Found %d ports: %v", len(ports), ports)
}

func TestGetPortsListType(t *testing.T) {
	ports, _ := GetPortsList()

	// Verify each port is a non-empty string
	for i, port := range ports {
		if port == "" {
			t.Errorf("Port %d is empty string", i)
		}
	}
}
