/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfigCmd(t *testing.T) {
	// Test that config command executes without error
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"config"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() error = %v", err)
	}
	// Note: The command prints to os.Stdout, not to the buffer
}

func TestConfigCmdDescription(t *testing.T) {
	if configCmd.Short == "" {
		t.Error("configCmd.Short should not be empty")
	}
	if configCmd.Long == "" {
		t.Error("configCmd.Long should not be empty")
	}
	if configCmd.Use != "config" {
		t.Errorf("configCmd.Use = %v, want config", configCmd.Use)
	}
}

func TestConfigCmdHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"config", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "config") {
		t.Errorf("help output should contain 'config', got: %s", output)
	}
}
