/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package transfer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewState(t *testing.T) {
	s := NewState()

	if s.Protocol != ProtocolZmodem {
		t.Errorf("NewState().Protocol = %d, want ProtocolZmodem", s.Protocol)
	}
	if s.InProgress {
		t.Error("NewState().InProgress should be false")
	}
	if s.Progress != 0 {
		t.Errorf("NewState().Progress = %f, want 0", s.Progress)
	}
	if s.InputFocus != 0 {
		t.Errorf("NewState().InputFocus = %d, want 0", s.InputFocus)
	}
}

func TestStateReset(t *testing.T) {
	s := NewState()
	s.InProgress = true
	s.Progress = 0.5
	s.StatusMsg = "Status"
	s.ErrorMsg = "Error"
	s.LocalPath = "/some/path"

	s.Reset()

	if s.InProgress {
		t.Error("Reset should set InProgress to false")
	}
	if s.Progress != 0 {
		t.Error("Reset should set Progress to 0")
	}
	if s.StatusMsg != "" {
		t.Error("Reset should clear StatusMsg")
	}
	if s.ErrorMsg != "" {
		t.Error("Reset should clear ErrorMsg")
	}
}

func TestGetProtocolName(t *testing.T) {
	tests := []struct {
		protocol Protocol
		want     string
	}{
		{ProtocolZmodem, "Zmodem"},
		{ProtocolXmodem, "Xmodem"},
		{ProtocolYmodem, "Ymodem"},
		{Protocol(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := NewState()
			s.Protocol = tt.protocol

			got := s.GetProtocolName()
			if got != tt.want {
				t.Errorf("GetProtocolName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCycleProtocol(t *testing.T) {
	s := NewState()

	// Start with Zmodem, cycle forward
	s.Protocol = ProtocolZmodem
	s.CycleProtocol(1)
	if s.Protocol != ProtocolXmodem {
		t.Error("CycleProtocol(1) from Zmodem should go to Xmodem")
	}

	s.CycleProtocol(1)
	if s.Protocol != ProtocolYmodem {
		t.Error("CycleProtocol(1) from Xmodem should go to Ymodem")
	}

	s.CycleProtocol(1)
	if s.Protocol != ProtocolZmodem {
		t.Error("CycleProtocol(1) from Ymodem should wrap to Zmodem")
	}

	// Cycle backward
	s.CycleProtocol(-1)
	if s.Protocol != ProtocolYmodem {
		t.Error("CycleProtocol(-1) from Zmodem should go to Ymodem")
	}
}

func TestFocusMethods(t *testing.T) {
	s := NewState()

	// FocusNext and FocusPrev are no-ops currently
	s.FocusNext()
	s.FocusPrev()
	// Just verify they don't panic
}

func TestProtocolConstants(t *testing.T) {
	// Verify protocol constants are distinct
	if ProtocolZmodem == ProtocolXmodem || ProtocolXmodem == ProtocolYmodem {
		t.Error("Protocol constants should be distinct")
	}

	if ProtocolZmodem != 0 {
		t.Error("ProtocolZmodem should be 0 (default)")
	}
}

func TestDirectionConstants(t *testing.T) {
	if DirectionSend == DirectionReceive {
		t.Error("Direction constants should be distinct")
	}

	if DirectionSend != 0 {
		t.Error("DirectionSend should be 0 (default)")
	}
}

func TestGetProtocolCommand(t *testing.T) {
	tests := []struct {
		direction Direction
		protocol  Protocol
		wantEmpty bool
	}{
		{DirectionSend, ProtocolZmodem, false},
		{DirectionReceive, ProtocolZmodem, false},
		{DirectionSend, ProtocolXmodem, false},
		{DirectionReceive, ProtocolXmodem, false},
		{DirectionSend, ProtocolYmodem, false},
		{DirectionReceive, ProtocolYmodem, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			cmd := GetProtocolCommand(tt.direction, tt.protocol)
			if tt.wantEmpty && cmd != "" {
				t.Errorf("GetProtocolCommand() = %q, want empty", cmd)
			}
			if !tt.wantEmpty && cmd == "" {
				t.Error("GetProtocolCommand() returned empty, want command")
			}
		})
	}
}

func TestDirectionLabel(t *testing.T) {
	if DirectionLabel(DirectionSend) != "send" {
		t.Errorf("DirectionLabel(Send) = %q, want 'send'", DirectionLabel(DirectionSend))
	}
	if DirectionLabel(DirectionReceive) != "receive" {
		t.Errorf("DirectionLabel(Receive) = %q, want 'receive'", DirectionLabel(DirectionReceive))
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			path:    "   ",
			wantErr: true,
		},
		{
			name:    "absolute path",
			path:    filepath.Join(os.TempDir(), "test"),
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "relative/path",
			wantErr: false,
		},
		{
			name:    "tilde expansion",
			path:    "~",
			wantErr: false,
		},
		{
			name:    "tilde with subpath",
			path:    "~/subdir",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandPath(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ExpandPath(%q) expected error, got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Errorf("ExpandPath(%q) unexpected error: %v", tt.path, err)
				}
				if result == "" {
					t.Errorf("ExpandPath(%q) returned empty string", tt.path)
				}
				// Result should be absolute
				if !filepath.IsAbs(result) {
					t.Errorf("ExpandPath(%q) = %q, expected absolute path", tt.path, result)
				}
			}
		})
	}
}

func TestExpandPathTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory")
	}

	result, err := ExpandPath("~")
	if err != nil {
		t.Fatalf("ExpandPath('~') error: %v", err)
	}

	if result != home {
		t.Errorf("ExpandPath('~') = %q, want %q", result, home)
	}
}

func TestExpandPathTildeSubdir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory")
	}

	result, err := ExpandPath("~/test")
	if err != nil {
		t.Fatalf("ExpandPath('~/test') error: %v", err)
	}

	expected := filepath.Join(home, "test")
	if result != expected {
		t.Errorf("ExpandPath('~/test') = %q, want %q", result, expected)
	}
}
