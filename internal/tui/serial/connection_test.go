/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package serial

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.BaudRate != DefaultBaudRate {
		t.Errorf("DefaultConfig().BaudRate = %d, want %d", cfg.BaudRate, DefaultBaudRate)
	}
	if cfg.DataBits != DefaultDataBits {
		t.Errorf("DefaultConfig().DataBits = %d, want %d", cfg.DataBits, DefaultDataBits)
	}
	if cfg.StopBits != DefaultStopBits {
		t.Errorf("DefaultConfig().StopBits = %d, want %d", cfg.StopBits, DefaultStopBits)
	}
	if cfg.Parity != DefaultParity {
		t.Errorf("DefaultConfig().Parity = %s, want %s", cfg.Parity, DefaultParity)
	}
	if cfg.LineEnding != DefaultLineEnding {
		t.Errorf("DefaultConfig().LineEnding = %s, want %s", cfg.LineEnding, DefaultLineEnding)
	}
	if cfg.Port != "" {
		t.Errorf("DefaultConfig().Port = %s, want empty string", cfg.Port)
	}
}

func TestNewConnection(t *testing.T) {
	tests := []struct {
		name         string
		port         string
		baudRate     int
		wantPort     string
		wantBaudRate int
	}{
		{
			name:         "empty port and default baud",
			port:         "",
			baudRate:     0,
			wantPort:     "",
			wantBaudRate: DefaultBaudRate,
		},
		{
			name:         "with port and custom baud",
			port:         "COM1",
			baudRate:     115200,
			wantPort:     "COM1",
			wantBaudRate: 115200,
		},
		{
			name:         "negative baud rate uses default",
			port:         "/dev/ttyUSB0",
			baudRate:     -1,
			wantPort:     "/dev/ttyUSB0",
			wantBaudRate: DefaultBaudRate,
		},
		{
			name:         "zero baud rate uses default",
			port:         "COM3",
			baudRate:     0,
			wantPort:     "COM3",
			wantBaudRate: DefaultBaudRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := NewConnection(tt.port, tt.baudRate)

			if conn == nil {
				t.Fatal("NewConnection returned nil")
			}
			if conn.Config.Port != tt.wantPort {
				t.Errorf("Config.Port = %s, want %s", conn.Config.Port, tt.wantPort)
			}
			if conn.Config.BaudRate != tt.wantBaudRate {
				t.Errorf("Config.BaudRate = %d, want %d", conn.Config.BaudRate, tt.wantBaudRate)
			}
			if conn.Connected {
				t.Error("NewConnection should not be connected")
			}
			if conn.Port != nil {
				t.Error("NewConnection should have nil Port")
			}
		})
	}
}

func TestConnectNoPort(t *testing.T) {
	conn := NewConnection("", 9600)

	err := conn.Connect()
	if err != ErrNoPortSelected {
		t.Errorf("Connect() error = %v, want %v", err, ErrNoPortSelected)
	}
	if conn.Connected {
		t.Error("Connection should not be connected after failed Connect")
	}
	if conn.LastError == "" {
		t.Error("LastError should be set after failed Connect")
	}
}

func TestConnectInvalidPort(t *testing.T) {
	conn := NewConnection("INVALID_PORT_THAT_DOES_NOT_EXIST", 9600)

	err := conn.Connect()
	if err == nil {
		conn.Disconnect()
		t.Error("Connect() should fail for invalid port")
	}
	if conn.Connected {
		t.Error("Connection should not be connected after failed Connect")
	}
}

func TestDisconnect(t *testing.T) {
	conn := NewConnection("COM1", 9600)

	// Disconnect without connection should not panic
	conn.Disconnect()

	if conn.Connected {
		t.Error("Connection should not be connected after Disconnect")
	}
	if conn.Port != nil {
		t.Error("Port should be nil after Disconnect")
	}
}

func TestWriteNoPort(t *testing.T) {
	conn := NewConnection("", 9600)

	n, err := conn.Write([]byte("test"))
	if err != ErrPortNotOpen {
		t.Errorf("Write() error = %v, want %v", err, ErrPortNotOpen)
	}
	if n != 0 {
		t.Errorf("Write() n = %d, want 0", n)
	}
}

func TestGetLineEndingBytes(t *testing.T) {
	tests := []struct {
		name       string
		lineEnding string
		want       string
	}{
		{
			name:       "None",
			lineEnding: "None",
			want:       "",
		},
		{
			name:       "NL",
			lineEnding: "NL",
			want:       "\n",
		},
		{
			name:       "CR",
			lineEnding: "CR",
			want:       "\r",
		},
		{
			name:       "CRLF",
			lineEnding: "CRLF",
			want:       "\r\n",
		},
		{
			name:       "default/unknown",
			lineEnding: "unknown",
			want:       "\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := NewConnection("", 9600)
			conn.Config.LineEnding = tt.lineEnding

			got := conn.GetLineEndingBytes()
			if got != tt.want {
				t.Errorf("GetLineEndingBytes() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetLineEndingDisplay(t *testing.T) {
	tests := []struct {
		name       string
		lineEnding string
		wantPart   string
	}{
		{
			name:       "None",
			lineEnding: "None",
			wantPart:   "No Line Ending",
		},
		{
			name:       "NL",
			lineEnding: "NL",
			wantPart:   "Newline",
		},
		{
			name:       "CR",
			lineEnding: "CR",
			wantPart:   "Carriage Return",
		},
		{
			name:       "CRLF",
			lineEnding: "CRLF",
			wantPart:   "NL & CR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := NewConnection("", 9600)
			conn.Config.LineEnding = tt.lineEnding

			got := conn.GetLineEndingDisplay()
			if got == "" {
				t.Error("GetLineEndingDisplay() returned empty string")
			}
		})
	}
}

func TestBaudRates(t *testing.T) {
	rates := BaudRates()

	if len(rates) == 0 {
		t.Fatal("BaudRates() returned empty slice")
	}

	// Check for common baud rates
	expectedRates := []int{9600, 115200}
	for _, expected := range expectedRates {
		found := false
		for _, rate := range rates {
			if rate == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BaudRates() should contain %d", expected)
		}
	}

	// Verify rates are positive
	for _, rate := range rates {
		if rate <= 0 {
			t.Errorf("BaudRates() contains non-positive rate: %d", rate)
		}
	}
}

func TestDataBitsOptions(t *testing.T) {
	options := DataBitsOptions()

	if len(options) == 0 {
		t.Fatal("DataBitsOptions() returned empty slice")
	}

	// Should contain 8 (most common)
	has8 := false
	for _, opt := range options {
		if opt == 8 {
			has8 = true
		}
		if opt < 5 || opt > 8 {
			t.Errorf("DataBitsOptions() contains invalid value: %d", opt)
		}
	}
	if !has8 {
		t.Error("DataBitsOptions() should contain 8")
	}
}

func TestStopBitsOptions(t *testing.T) {
	options := StopBitsOptions()

	if len(options) == 0 {
		t.Fatal("StopBitsOptions() returned empty slice")
	}

	// Should contain 1 and 2
	has1 := false
	has2 := false
	for _, opt := range options {
		if opt == 1 {
			has1 = true
		}
		if opt == 2 {
			has2 = true
		}
	}
	if !has1 || !has2 {
		t.Error("StopBitsOptions() should contain 1 and 2")
	}
}

func TestLineEndingOptions(t *testing.T) {
	options := LineEndingOptions()

	if len(options) == 0 {
		t.Fatal("LineEndingOptions() returned empty slice")
	}

	expectedOptions := []string{"None", "NL", "CR", "CRLF"}
	for _, expected := range expectedOptions {
		found := false
		for _, opt := range options {
			if opt == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("LineEndingOptions() should contain %s", expected)
		}
	}
}

func TestConstants(t *testing.T) {
	if DefaultBaudRate != 9600 {
		t.Errorf("DefaultBaudRate = %d, want 9600", DefaultBaudRate)
	}
	if DefaultDataBits != 8 {
		t.Errorf("DefaultDataBits = %d, want 8", DefaultDataBits)
	}
	if DefaultStopBits != 1 {
		t.Errorf("DefaultStopBits = %d, want 1", DefaultStopBits)
	}
	if DefaultParity != "None" {
		t.Errorf("DefaultParity = %s, want None", DefaultParity)
	}
	if DefaultLineEnding != "CRLF" {
		t.Errorf("DefaultLineEnding = %s, want CRLF", DefaultLineEnding)
	}
}

func TestErrors(t *testing.T) {
	if ErrNoPortSelected == nil {
		t.Error("ErrNoPortSelected should not be nil")
	}
	if ErrPortNotOpen == nil {
		t.Error("ErrPortNotOpen should not be nil")
	}

	if ErrNoPortSelected.Error() == "" {
		t.Error("ErrNoPortSelected.Error() should not be empty")
	}
	if ErrPortNotOpen.Error() == "" {
		t.Error("ErrPortNotOpen.Error() should not be empty")
	}
}
