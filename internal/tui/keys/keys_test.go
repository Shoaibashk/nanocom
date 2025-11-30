/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package keys

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestDefaultKeys(t *testing.T) {
	keys := DefaultKeys()

	tests := []struct {
		name    string
		binding key.Binding
		wantKey string
	}{
		{
			name:    "Quit binding has q",
			binding: keys.Quit,
			wantKey: "q",
		},
		{
			name:    "Menu binding has ctrl+a",
			binding: keys.Menu,
			wantKey: "ctrl+a",
		},
		{
			name:    "Enter binding has enter",
			binding: keys.Enter,
			wantKey: "enter",
		},
		{
			name:    "Up binding has up",
			binding: keys.Up,
			wantKey: "up",
		},
		{
			name:    "Down binding has down",
			binding: keys.Down,
			wantKey: "down",
		},
		{
			name:    "Escape binding has esc",
			binding: keys.Escape,
			wantKey: "esc",
		},
		{
			name:    "CtrlA binding has ctrl+a",
			binding: keys.CtrlA,
			wantKey: "ctrl+a",
		},
		{
			name:    "Help binding has z",
			binding: keys.Help,
			wantKey: "z",
		},
		{
			name:    "Connect binding has c",
			binding: keys.Connect,
			wantKey: "c",
		},
		{
			name:    "Disconnect binding has d",
			binding: keys.Disconnect,
			wantKey: "d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, k := range tt.binding.Keys() {
				if k == tt.wantKey {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("binding does not contain key %q, has: %v", tt.wantKey, tt.binding.Keys())
			}
		})
	}
}

func TestDefaultKeysEnabled(t *testing.T) {
	keys := DefaultKeys()

	bindings := []key.Binding{
		keys.Quit,
		keys.Menu,
		keys.Enter,
		keys.Up,
		keys.Down,
		keys.Escape,
		keys.CtrlA,
		keys.Help,
		keys.Connect,
		keys.Disconnect,
	}

	for i, binding := range bindings {
		if !binding.Enabled() {
			t.Errorf("binding %d should be enabled by default", i)
		}
	}
}

func TestDefaultKeysHasHelp(t *testing.T) {
	keys := DefaultKeys()

	bindings := []struct {
		name    string
		binding key.Binding
	}{
		{"Quit", keys.Quit},
		{"Menu", keys.Menu},
		{"Enter", keys.Enter},
		{"Up", keys.Up},
		{"Down", keys.Down},
		{"Escape", keys.Escape},
		{"CtrlA", keys.CtrlA},
		{"Help", keys.Help},
		{"Connect", keys.Connect},
		{"Disconnect", keys.Disconnect},
	}

	for _, b := range bindings {
		t.Run(b.name, func(t *testing.T) {
			help := b.binding.Help()
			if help.Key == "" {
				t.Errorf("%s binding should have help key", b.name)
			}
			if help.Desc == "" {
				t.Errorf("%s binding should have help description", b.name)
			}
		})
	}
}

func TestMapStruct(t *testing.T) {
	m := Map{}

	// Test that empty Map has zero values
	if m.Quit.Enabled() {
		t.Error("empty Map.Quit should not be enabled")
	}

	// Test that DefaultKeys returns proper Map
	defaultMap := DefaultKeys()
	if !defaultMap.Quit.Enabled() {
		t.Error("DefaultKeys().Quit should be enabled")
	}
}

func TestMultipleKeyBindings(t *testing.T) {
	keys := DefaultKeys()

	// Quit should have multiple keys (q and ctrl+c)
	quitKeys := keys.Quit.Keys()
	if len(quitKeys) < 2 {
		t.Errorf("Quit binding should have multiple keys, got: %v", quitKeys)
	}

	hasQ := false
	hasCtrlC := false
	for _, k := range quitKeys {
		if k == "q" {
			hasQ = true
		}
		if k == "ctrl+c" {
			hasCtrlC = true
		}
	}

	if !hasQ {
		t.Error("Quit binding should include 'q'")
	}
	if !hasCtrlC {
		t.Error("Quit binding should include 'ctrl+c'")
	}
}

func TestUpDownAlternatives(t *testing.T) {
	keys := DefaultKeys()

	// Up should have arrow up and k (vim-style)
	upKeys := keys.Up.Keys()
	hasUp := false
	hasK := false
	for _, k := range upKeys {
		if k == "up" {
			hasUp = true
		}
		if k == "k" {
			hasK = true
		}
	}
	if !hasUp || !hasK {
		t.Errorf("Up binding should have 'up' and 'k', got: %v", upKeys)
	}

	// Down should have arrow down and j (vim-style)
	downKeys := keys.Down.Keys()
	hasDown := false
	hasJ := false
	for _, k := range downKeys {
		if k == "down" {
			hasDown = true
		}
		if k == "j" {
			hasJ = true
		}
	}
	if !hasDown || !hasJ {
		t.Errorf("Down binding should have 'down' and 'j', got: %v", downKeys)
	}
}
