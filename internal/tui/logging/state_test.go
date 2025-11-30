/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewState(t *testing.T) {
	s := NewState()

	if s.IsEnabled() {
		t.Error("NewState should return disabled state")
	}
	if s.GetFilePath() != "" {
		t.Error("NewState should have empty file path")
	}

	lines, bytes := s.GetStats()
	if lines != 0 || bytes != 0 {
		t.Errorf("NewState stats should be 0, got lines=%d, bytes=%d", lines, bytes)
	}
}

func TestStartStop(t *testing.T) {
	s := NewState()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Start logging
	err := s.Start(logPath)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if !s.IsEnabled() {
		t.Error("State should be enabled after Start")
	}
	if s.GetFilePath() != logPath {
		t.Errorf("GetFilePath() = %v, want %v", s.GetFilePath(), logPath)
	}

	// Stop logging
	err = s.Stop()
	if err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	if s.IsEnabled() {
		t.Error("State should be disabled after Stop")
	}

	// Verify file was created and contains header/footer
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "nanocom capture started") {
		t.Error("Log file should contain header")
	}
	if !strings.Contains(string(content), "nanocom capture ended") {
		t.Error("Log file should contain footer")
	}
}

func TestStartAlreadyStarted(t *testing.T) {
	s := NewState()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	err := s.Start(logPath)
	if err != nil {
		t.Fatalf("First Start() error = %v", err)
	}
	defer s.Stop()

	// Second start should be no-op
	err = s.Start(logPath)
	if err != nil {
		t.Errorf("Second Start() should not error, got = %v", err)
	}
}

func TestStopNotStarted(t *testing.T) {
	s := NewState()

	// Stop on not-started should be no-op
	err := s.Stop()
	if err != nil {
		t.Errorf("Stop() on not-started state should not error, got = %v", err)
	}
}

func TestWrite(t *testing.T) {
	s := NewState()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	err := s.Start(logPath)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	testData := "Hello, World!"
	err = s.Write(testData)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	lines, bytes := s.GetStats()
	if bytes != int64(len(testData)) {
		t.Errorf("bytes = %d, want %d", bytes, len(testData))
	}
	if lines != 0 {
		t.Errorf("lines = %d, want 0 (no newlines written)", lines)
	}

	err = s.Stop()
	if err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Verify content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), testData) {
		t.Errorf("Log file should contain %q", testData)
	}
}

func TestWriteLine(t *testing.T) {
	s := NewState()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	err := s.Start(logPath)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	err = s.WriteLine("Line 1")
	if err != nil {
		t.Errorf("WriteLine() error = %v", err)
	}

	err = s.WriteLine("Line 2")
	if err != nil {
		t.Errorf("WriteLine() error = %v", err)
	}

	lines, _ := s.GetStats()
	if lines != 2 {
		t.Errorf("lines = %d, want 2", lines)
	}

	err = s.Stop()
	if err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Verify content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Line 1") {
		t.Error("Log file should contain 'Line 1'")
	}
	if !strings.Contains(string(content), "Line 2") {
		t.Error("Log file should contain 'Line 2'")
	}
}

func TestWriteWhenDisabled(t *testing.T) {
	s := NewState()

	// Write when not started should be no-op
	err := s.Write("test")
	if err != nil {
		t.Errorf("Write() when disabled should not error, got = %v", err)
	}

	err = s.WriteLine("test")
	if err != nil {
		t.Errorf("WriteLine() when disabled should not error, got = %v", err)
	}
}

func TestToggle(t *testing.T) {
	s := NewState()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Toggle on
	err := s.Toggle(logPath)
	if err != nil {
		t.Fatalf("Toggle() on error = %v", err)
	}
	if !s.IsEnabled() {
		t.Error("State should be enabled after first Toggle")
	}

	// Toggle off
	err = s.Toggle(logPath)
	if err != nil {
		t.Fatalf("Toggle() off error = %v", err)
	}
	if s.IsEnabled() {
		t.Error("State should be disabled after second Toggle")
	}
}

func TestStartWithDefaultPath(t *testing.T) {
	s := NewState()

	// Start with empty path should use default
	err := s.Start("")
	if err != nil {
		t.Fatalf("Start() with empty path error = %v", err)
	}

	path := s.GetFilePath()
	if path == "" {
		t.Error("GetFilePath() should not be empty after Start with default")
	}

	if !strings.Contains(path, "nanocom_capture.log") {
		t.Errorf("Default path should contain 'nanocom_capture.log', got %s", path)
	}

	err = s.Stop()
	if err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Clean up the file
	os.Remove(path)
}

func TestStartWithInvalidPath(t *testing.T) {
	s := NewState()

	// Try to start with invalid path (directory that doesn't exist and can't be created)
	err := s.Start("/nonexistent/directory/that/cannot/exist/file.log")
	if err == nil {
		s.Stop()
		t.Error("Start() with invalid path should return error")
	}

	if s.IsEnabled() {
		t.Error("State should not be enabled after failed Start")
	}
}

func TestGetStatsAfterMultipleWrites(t *testing.T) {
	s := NewState()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	err := s.Start(logPath)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Write multiple lines
	for i := 0; i < 5; i++ {
		s.WriteLine("Test line")
	}

	lines, bytes := s.GetStats()
	if lines != 5 {
		t.Errorf("lines = %d, want 5", lines)
	}
	expectedBytes := int64(len("Test line\n") * 5)
	if bytes != expectedBytes {
		t.Errorf("bytes = %d, want %d", bytes, expectedBytes)
	}

	s.Stop()
}
