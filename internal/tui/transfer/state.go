/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package transfer

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Protocol represents the file transfer protocol
type Protocol int

const (
	ProtocolZmodem Protocol = iota
	ProtocolXmodem
	ProtocolYmodem
)

// Direction represents whether we're sending or receiving
type Direction int

const (
	DirectionSend Direction = iota
	DirectionReceive
)

// State holds the current file transfer state
type State struct {
	Direction  Direction
	Protocol   Protocol
	LocalPath  string
	InProgress bool
	Progress   float64
	StatusMsg  string
	ErrorMsg   string
	FileInput  textinput.Model
	InputFocus int // 0=file
}

// PrimaryColor for input styling
var primaryColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

// NewState creates a new transfer state
func NewState() State {
	fileInput := textinput.New()
	fileInput.Placeholder = "Local file path..."
	fileInput.CharLimit = 256
	fileInput.Width = 50
	fileInput.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)

	return State{
		Protocol:   ProtocolZmodem,
		FileInput:  fileInput,
		InputFocus: 0,
	}
}

// FocusNext moves focus to next input field
func (s *State) FocusNext() {
	// Only file input, no-op
}

// FocusPrev moves focus to previous input field
func (s *State) FocusPrev() {
	// Only file input, no-op
}

func (s *State) updateFocus() {
	s.FileInput.Focus()
}

// Reset clears the transfer state
func (s *State) Reset() {
	s.InProgress = false
	s.Progress = 0
	s.StatusMsg = ""
	s.ErrorMsg = ""
	s.FileInput.SetValue("")
	s.InputFocus = 0
	s.FileInput.Focus()
}

// GetProtocolName returns the display name of the protocol
func (s *State) GetProtocolName() string {
	switch s.Protocol {
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
func (s *State) CycleProtocol(delta int) {
	protocols := []Protocol{ProtocolZmodem, ProtocolXmodem, ProtocolYmodem}
	currentIdx := 0
	for i, p := range protocols {
		if p == s.Protocol {
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
	s.Protocol = protocols[currentIdx]
}

// GetProtocolCommand returns the command for the given direction and protocol
func GetProtocolCommand(direction Direction, protocol Protocol) string {
	switch protocol {
	case ProtocolZmodem:
		if direction == DirectionSend {
			return "rz -y"
		}
		return "sz"
	case ProtocolXmodem:
		if direction == DirectionSend {
			return "rx"
		}
		return "sx"
	case ProtocolYmodem:
		if direction == DirectionSend {
			return "rb -y"
		}
		return "sb"
	default:
		return ""
	}
}

// DirectionLabel returns a human-readable label for the direction
func DirectionLabel(direction Direction) string {
	if direction == DirectionSend {
		return "send"
	}
	return "receive"
}
