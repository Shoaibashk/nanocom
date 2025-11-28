/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// TransferProtocol represents the file transfer protocol
type TransferProtocol int

const (
	ProtocolZmodem TransferProtocol = iota
	ProtocolXmodem
	ProtocolYmodem
)

// TransferDirection represents whether we're sending or receiving
type TransferDirection int

const (
	TransferSend TransferDirection = iota
	TransferReceive
)

// TransferState holds the current file transfer state
type TransferState struct {
	direction  TransferDirection
	protocol   TransferProtocol
	localPath  string
	inProgress bool
	progress   float64
	statusMsg  string
	errorMsg   string
	fileInput  textinput.Model
	inputFocus int // 0=file
}

// NewTransferState creates a new transfer state
func NewTransferState() TransferState {
	fileInput := textinput.New()
	fileInput.Placeholder = "Local file path..."
	fileInput.CharLimit = 256
	fileInput.Width = 50
	fileInput.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)

	return TransferState{
		protocol:   ProtocolZmodem,
		fileInput:  fileInput,
		inputFocus: 0,
	}
}

// FocusNext moves focus to next input field
func (t *TransferState) FocusNext() {
	// Only file input, no-op
}

// FocusPrev moves focus to previous input field
func (t *TransferState) FocusPrev() {
	// Only file input, no-op
}

func (t *TransferState) updateFocus() {
	t.fileInput.Focus()
}

// Reset clears the transfer state
func (t *TransferState) Reset() {
	t.inProgress = false
	t.progress = 0
	t.statusMsg = ""
	t.errorMsg = ""
	t.fileInput.SetValue("")
	t.inputFocus = 0
	t.fileInput.Focus()
}

// GetProtocolName returns the display name of the protocol
func (t *TransferState) GetProtocolName() string {
	switch t.protocol {
	case ProtocolZmodem:
		return "Zmodem"
	case ProtocolXmodem:
		return "Xmodem"
	case ProtocolYmodem:
		return "Ymodem"
	default:
		return "Unknown"
	}
}

// CycleProtocol cycles through available protocols
func (t *TransferState) CycleProtocol(delta int) {
	protocols := []TransferProtocol{ProtocolZmodem, ProtocolXmodem, ProtocolYmodem}
	currentIdx := 0
	for i, p := range protocols {
		if p == t.protocol {
			currentIdx = i
			break
		}
	}
	currentIdx += delta
	if currentIdx < 0 {
		currentIdx = len(protocols) - 1
	} else if currentIdx >= len(protocols) {
		currentIdx = 0
	}
	t.protocol = protocols[currentIdx]
}

// LoggingState holds the logging state
type LoggingState struct {
	enabled   bool
	filePath  string
	file      *os.File
	lineCount int
	byteCount int64
}

// NewLoggingState creates a new logging state
func NewLoggingState() LoggingState {
	return LoggingState{
		enabled: false,
	}
}

// Toggle toggles logging on/off
func (l *LoggingState) Toggle(defaultPath string) error {
	if l.enabled {
		return l.Stop()
	}
	return l.Start(defaultPath)
}

// Start starts logging to the specified file
func (l *LoggingState) Start(filePath string) error {
	if l.enabled {
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

	l.file = file
	l.filePath = filePath
	l.enabled = true
	l.lineCount = 0
	l.byteCount = 0

	return nil
}

// Stop stops logging
func (l *LoggingState) Stop() error {
	if !l.enabled {
		return nil
	}

	if l.file != nil {
		// Write footer
		footer := fmt.Sprintf("\n=== nanocom capture ended (lines: %d, bytes: %d) ===\n", l.lineCount, l.byteCount)
		l.file.WriteString(footer)
		l.file.Close()
		l.file = nil
	}

	l.enabled = false
	return nil
}

// Write writes data to the log file
func (l *LoggingState) Write(data string) error {
	if !l.enabled || l.file == nil {
		return nil
	}

	n, err := io.WriteString(l.file, data)
	if err != nil {
		return err
	}

	l.byteCount += int64(n)
	l.lineCount += strings.Count(data, "\n")

	return nil
}

// WriteLine writes a line to the log file
func (l *LoggingState) WriteLine(line string) error {
	return l.Write(line + "\n")
}

// IsEnabled returns whether logging is enabled
func (l *LoggingState) IsEnabled() bool {
	return l.enabled
}

// GetFilePath returns the current log file path
func (l *LoggingState) GetFilePath() string {
	return l.filePath
}

// GetStats returns logging statistics
func (l *LoggingState) GetStats() (lineCount int, byteCount int64) {
	return l.lineCount, l.byteCount
}
