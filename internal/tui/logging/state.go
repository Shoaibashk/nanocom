/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// State holds the logging state
type State struct {
	enabled   bool
	filePath  string
	file      *os.File
	lineCount int
	byteCount int64
}

// NewState creates a new logging state
func NewState() State {
	return State{
		enabled: false,
	}
}

// Toggle toggles logging on/off
func (s *State) Toggle(defaultPath string) error {
	if s.enabled {
		return s.Stop()
	}
	return s.Start(defaultPath)
}

// Start starts logging to the specified file
func (s *State) Start(filePath string) error {
	if s.enabled {
		return nil // Already logging
	}

	// Use default path if not specified
	if filePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		filePath = filepath.Join(homeDir, "nanocom_capture.log")
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Write header
	header := "\n=== nanocom capture started ===\n"
	file.WriteString(header)

	s.file = file
	s.filePath = filePath
	s.enabled = true
	s.lineCount = 0
	s.byteCount = 0

	return nil
}

// Stop stops logging
func (s *State) Stop() error {
	if !s.enabled {
		return nil
	}

	if s.file != nil {
		// Write footer
		footer := fmt.Sprintf("\n=== nanocom capture ended (lines: %d, bytes: %d) ===\n", s.lineCount, s.byteCount)
		s.file.WriteString(footer)
		s.file.Close()
		s.file = nil
	}

	s.enabled = false
	return nil
}

// Write writes data to the log file
func (s *State) Write(data string) error {
	if !s.enabled || s.file == nil {
		return nil
	}

	n, err := io.WriteString(s.file, data)
	if err != nil {
		return err
	}

	s.byteCount += int64(n)
	s.lineCount += strings.Count(data, "\n")

	return nil
}

// WriteLine writes a line to the log file
func (s *State) WriteLine(line string) error {
	return s.Write(line + "\n")
}

// IsEnabled returns whether logging is enabled
func (s *State) IsEnabled() bool {
	return s.enabled
}

// GetFilePath returns the current log file path
func (s *State) GetFilePath() string {
	return s.filePath
}

// GetStats returns logging statistics
func (s *State) GetStats() (lineCount int, byteCount int64) {
	return s.lineCount, s.byteCount
}
