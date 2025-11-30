/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestListCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() error = %v", err)
	}
	// Note: The list command prints to os.Stdout, not to the buffer
	// Just verify it runs without error
}

func TestListCmdDescription(t *testing.T) {
	if listCmd.Short == "" {
		t.Error("listCmd.Short should not be empty")
	}
	if listCmd.Long == "" {
		t.Error("listCmd.Long should not be empty")
	}
	if listCmd.Use != "list" {
		t.Errorf("listCmd.Use = %v, want list", listCmd.Use)
	}
}

func TestListCmdHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("rootCmd.Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "list") {
		t.Errorf("help output should contain 'list', got: %s", output)
	}
	if !strings.Contains(output, "serial ports") {
		t.Errorf("help output should contain 'serial ports', got: %s", output)
	}
}
