/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "help flag",
			args: []string{"--help"},
		},
		{
			name: "version flag short",
			args: []string{"-h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			// Reset flags for each test
			baudrate = 9600
			port = ""

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("rootCmd.Execute() error = %v", err)
			}
		})
	}
}

func TestRootCmdFlags(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantBaudrate int
		wantPort     string
	}{
		{
			name:         "default values",
			args:         []string{"--help"},
			wantBaudrate: 9600,
			wantPort:     "",
		},
		{
			name:         "custom baudrate",
			args:         []string{"-b", "115200", "--help"},
			wantBaudrate: 115200,
			wantPort:     "",
		},
		{
			name:         "custom port",
			args:         []string{"-p", "COM3", "--help"},
			wantBaudrate: 9600,
			wantPort:     "COM3",
		},
		{
			name:         "both flags",
			args:         []string{"-p", "/dev/ttyUSB0", "-b", "57600", "--help"},
			wantBaudrate: 57600,
			wantPort:     "/dev/ttyUSB0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			// Reset flags for each test
			baudrate = 9600
			port = ""

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("rootCmd.Execute() error = %v", err)
			}

			if baudrate != tt.wantBaudrate {
				t.Errorf("baudrate = %v, want %v", baudrate, tt.wantBaudrate)
			}
			if port != tt.wantPort {
				t.Errorf("port = %v, want %v", port, tt.wantPort)
			}
		})
	}
}

func TestRootCmdDescription(t *testing.T) {
	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}
	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
	if rootCmd.Use != "nanocom" {
		t.Errorf("rootCmd.Use = %v, want nanocom", rootCmd.Use)
	}
}
