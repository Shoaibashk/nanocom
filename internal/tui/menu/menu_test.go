/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package menu

import (
	"testing"
)

func TestItem(t *testing.T) {
	item := Item{
		Key:   "T",
		Label: "Test Item",
	}

	if item.Key != "T" {
		t.Errorf("Item.Key = %v, want T", item.Key)
	}
	if item.Label != "Test Item" {
		t.Errorf("Item.Label = %v, want 'Test Item'", item.Label)
	}
}

func TestSerialSetupConstants(t *testing.T) {
	// Verify constants are defined correctly
	if SerialSetupPort != 0 {
		t.Errorf("SerialSetupPort = %d, want 0", SerialSetupPort)
	}
	if SerialSetupBaudRate != 1 {
		t.Errorf("SerialSetupBaudRate = %d, want 1", SerialSetupBaudRate)
	}
	if SerialSetupDataBits != 2 {
		t.Errorf("SerialSetupDataBits = %d, want 2", SerialSetupDataBits)
	}
	if SerialSetupStopBits != 3 {
		t.Errorf("SerialSetupStopBits = %d, want 3", SerialSetupStopBits)
	}
	if SerialSetupLineEnding != 4 {
		t.Errorf("SerialSetupLineEnding = %d, want 4", SerialSetupLineEnding)
	}
	if SerialSetupConnect != 5 {
		t.Errorf("SerialSetupConnect = %d, want 5", SerialSetupConnect)
	}
	if SerialSetupMenuItemCount != 6 {
		t.Errorf("SerialSetupMenuItemCount = %d, want 6", SerialSetupMenuItemCount)
	}
}

func TestDefaultMainMenuItems(t *testing.T) {
	items := DefaultMainMenuItems()

	if len(items) == 0 {
		t.Fatal("DefaultMainMenuItems should return non-empty slice")
	}

	// Check for essential menu items
	expectedKeys := map[string]bool{
		"O": false, // Serial port settings
		"S": false, // Send files
		"R": false, // Receive files
		"C": false, // Clear screen
		"L": false, // Log
		"Z": false, // Help
		"X": false, // Exit
		"Q": false, // Quit
	}

	for _, item := range items {
		if item.Key == "" {
			t.Error("Menu item should have a key")
		}
		if item.Label == "" {
			t.Error("Menu item should have a label")
		}
		expectedKeys[item.Key] = true
	}

	for key, found := range expectedKeys {
		if !found {
			t.Errorf("Expected menu item with key %q not found", key)
		}
	}
}

func TestDefaultHelpMenuItems(t *testing.T) {
	items := DefaultHelpMenuItems()

	if len(items) == 0 {
		t.Fatal("DefaultHelpMenuItems should return non-empty slice")
	}

	// All items should have key and label
	for i, item := range items {
		if item.Key == "" {
			t.Errorf("Help menu item %d should have a key", i)
		}
		if item.Label == "" {
			t.Errorf("Help menu item %d should have a label", i)
		}
	}

	// Check that help menu includes Ctrl+A commands
	hasCtrlA := false
	for _, item := range items {
		if len(item.Key) > 0 && item.Key[0] == '^' {
			hasCtrlA = true
			break
		}
	}
	if !hasCtrlA {
		t.Error("Help menu should include Ctrl+A commands (^A prefix)")
	}
}

func TestMainMenuItemsAreUnique(t *testing.T) {
	items := DefaultMainMenuItems()

	keys := make(map[string]bool)
	for _, item := range items {
		if keys[item.Key] {
			t.Errorf("Duplicate menu key: %s", item.Key)
		}
		keys[item.Key] = true
	}
}

func TestHelpMenuItemsAreUnique(t *testing.T) {
	items := DefaultHelpMenuItems()

	keys := make(map[string]bool)
	for _, item := range items {
		if keys[item.Key] {
			t.Errorf("Duplicate help menu key: %s", item.Key)
		}
		keys[item.Key] = true
	}
}

func TestMainMenuItemOrder(t *testing.T) {
	items := DefaultMainMenuItems()

	// First item should be serial port settings (O)
	if items[0].Key != "O" {
		t.Errorf("First menu item should be 'O', got %s", items[0].Key)
	}

	// Last item should be quit (Q)
	lastItem := items[len(items)-1]
	if lastItem.Key != "Q" {
		t.Errorf("Last menu item should be 'Q', got %s", lastItem.Key)
	}
}

func TestHelpMenuHasQuit(t *testing.T) {
	items := DefaultHelpMenuItems()

	hasQuit := false
	for _, item := range items {
		if item.Key == "^A Q" {
			hasQuit = true
			break
		}
	}

	if !hasQuit {
		t.Error("Help menu should include quit command (^A Q)")
	}
}
