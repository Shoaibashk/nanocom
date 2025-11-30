/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package update

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/shoaibashk/nanocom/internal/tui/model"
	"github.com/shoaibashk/nanocom/internal/tui/state"
	"github.com/shoaibashk/nanocom/internal/tui/transfer"
)

func TestUpdateWindowSizeMsg(t *testing.T) {
	m := model.New("", 0)
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	newModel, _ := Update(m, msg)

	if newModel.Width != 100 {
		t.Errorf("Width = %d, want 100", newModel.Width)
	}
	if newModel.Height != 50 {
		t.Errorf("Height = %d, want 50", newModel.Height)
	}
}

func TestUpdatePortsListMsg(t *testing.T) {
	m := model.New("", 0)
	ports := []string{"COM1", "COM2", "COM3"}
	msg := state.PortsListMsg(ports)

	newModel, _ := Update(m, msg)

	if len(newModel.AvailablePorts) != 3 {
		t.Errorf("AvailablePorts length = %d, want 3", len(newModel.AvailablePorts))
	}
}

func TestUpdateSerialDataMsg(t *testing.T) {
	m := model.New("", 0)
	msg := state.SerialDataMsg("Hello\n")

	newModel, _ := Update(m, msg)

	// Data should be processed (may or may not appear depending on buffer state)
	_ = newModel
}

func TestUpdateTickMsg(t *testing.T) {
	m := model.New("", 0)
	m.TransferPending = true

	msg := state.TickMsg{}

	newModel, _ := Update(m, msg)

	// Spinner should advance
	_ = newModel
}

func TestUpdateAutoConnectMsg(t *testing.T) {
	m := model.New("COM1", 9600)
	m.Serial.Config.Port = "COM1"
	msg := state.AutoConnectMsg{}

	newModel, _ := Update(m, msg)

	// Should attempt to connect (will fail without real port)
	_ = newModel
}

func TestHandleCtrlACommands(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantView state.ViewState
	}{
		{
			name:     "z opens help",
			key:      "z",
			wantView: state.ViewHelp,
		},
		{
			name:     "c clears screen",
			key:      "c",
			wantView: state.ViewTerminal,
		},
		{
			name:     "o opens main menu",
			key:      "o",
			wantView: state.ViewMainMenu,
		},
		{
			name:     "s opens send file",
			key:      "s",
			wantView: state.ViewSendFile,
		},
		{
			name:     "r opens receive file",
			key:      "r",
			wantView: state.ViewReceiveFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.New("", 0)
			m.CtrlAPressed = true

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newModel, _ := Update(m, msg)

			if newModel.CurrentView != tt.wantView {
				t.Errorf("CurrentView = %v, want %v", newModel.CurrentView, tt.wantView)
			}
			if newModel.CtrlAPressed {
				t.Error("CtrlAPressed should be false after command")
			}
		})
	}
}

func TestHandleCtrlAPressed(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewTerminal

	msg := tea.KeyMsg{Type: tea.KeyCtrlA}

	newModel, _ := Update(m, msg)

	if !newModel.CtrlAPressed {
		t.Error("CtrlAPressed should be true after Ctrl+A")
	}
}

func TestHandleMenuNavigation(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewMainMenu
	m.MenuIndex = 0

	// Test down
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := Update(m, msg)

	if newModel.MenuIndex != 1 {
		t.Errorf("MenuIndex = %d, want 1 after down", newModel.MenuIndex)
	}

	// Test up
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = Update(newModel, msg)

	if newModel.MenuIndex != 0 {
		t.Errorf("MenuIndex = %d, want 0 after up", newModel.MenuIndex)
	}
}

func TestHandleEscape(t *testing.T) {
	tests := []struct {
		name      string
		startView state.ViewState
		wantView  state.ViewState
	}{
		{
			name:      "escape from main menu",
			startView: state.ViewMainMenu,
			wantView:  state.ViewTerminal,
		},
		{
			name:      "escape from serial setup",
			startView: state.ViewSerialSetup,
			wantView:  state.ViewTerminal,
		},
		{
			name:      "escape from help",
			startView: state.ViewHelp,
			wantView:  state.ViewTerminal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.New("", 0)
			m.CurrentView = tt.startView

			msg := tea.KeyMsg{Type: tea.KeyEscape}
			newModel, _ := Update(m, msg)

			if newModel.CurrentView != tt.wantView {
				t.Errorf("CurrentView = %v, want %v", newModel.CurrentView, tt.wantView)
			}
		})
	}
}

func TestTransferResultMsgSuccess(t *testing.T) {
	m := model.New("", 0)
	m.TransferPending = true

	msg := TransferResultMsg{
		Direction: transfer.DirectionSend,
		Path:      "/path/to/file.txt",
		Err:       nil,
	}

	newModel, _ := Update(m, msg)

	if newModel.TransferPending {
		t.Error("TransferPending should be false after result")
	}
}

func TestTransferResultMsgError(t *testing.T) {
	m := model.New("", 0)
	m.TransferPending = true

	msg := TransferResultMsg{
		Direction: transfer.DirectionSend,
		Path:      "/path/to/file.txt",
		Err:       &testError{},
	}

	newModel, _ := Update(m, msg)

	if newModel.TransferPending {
		t.Error("TransferPending should be false after error result")
	}
}

type testError struct{}

func (e *testError) Error() string {
	return "test error"
}

func TestProcessSerialData(t *testing.T) {
	m := model.New("", 0)

	// Process data with newlines
	msg := state.SerialDataMsg("Line 1\nLine 2\n")
	newModel, _ := Update(m, msg)

	// Buffer should contain the lines
	_ = newModel
}

func TestHandlePortListNavigation(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewPortList
	m.AvailablePorts = []string{"COM1", "COM2", "COM3"}
	m.PortIndex = 0

	// Navigate down
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := Update(m, msg)

	if newModel.PortIndex != 1 {
		t.Errorf("PortIndex = %d, want 1", newModel.PortIndex)
	}

	// Navigate up
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = Update(newModel, msg)

	if newModel.PortIndex != 0 {
		t.Errorf("PortIndex = %d, want 0", newModel.PortIndex)
	}
}

func TestHandlePortListSelect(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewPortList
	m.AvailablePorts = []string{"COM1", "COM2", "COM3"}
	m.PortIndex = 1

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := Update(m, msg)

	if newModel.Serial.Config.Port != "COM2" {
		t.Errorf("Port = %q, want 'COM2'", newModel.Serial.Config.Port)
	}
	if newModel.CurrentView != state.ViewSerialSetup {
		t.Errorf("CurrentView = %v, want ViewSerialSetup", newModel.CurrentView)
	}
}

func TestHandleSendFileEscape(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewSendFile

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	newModel, _ := Update(m, msg)

	if newModel.CurrentView != state.ViewTerminal {
		t.Errorf("CurrentView = %v, want ViewTerminal", newModel.CurrentView)
	}
}

func TestHandleReceiveFileEscape(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewReceiveFile

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	newModel, _ := Update(m, msg)

	if newModel.CurrentView != state.ViewTerminal {
		t.Errorf("CurrentView = %v, want ViewTerminal", newModel.CurrentView)
	}
}

func TestHandleSerialSetupNavigation(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewSerialSetup
	m.MenuIndex = 0

	// Navigate down
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := Update(m, msg)

	if newModel.MenuIndex != 1 {
		t.Errorf("MenuIndex = %d, want 1", newModel.MenuIndex)
	}
}

func TestHandleSerialSetupAdjust(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewSerialSetup
	m.MenuIndex = 1 // Baud rate

	initialBaud := m.Serial.Config.BaudRate

	// Adjust right
	msg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ := Update(m, msg)

	if newModel.Serial.Config.BaudRate == initialBaud {
		t.Error("BaudRate should change with right arrow")
	}
}

func TestVimStyleNavigation(t *testing.T) {
	m := model.New("", 0)
	m.CurrentView = state.ViewMainMenu
	m.MenuIndex = 0

	// Test j for down
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	newModel, _ := Update(m, msg)

	if newModel.MenuIndex != 1 {
		t.Errorf("MenuIndex = %d, want 1 after 'j'", newModel.MenuIndex)
	}

	// Test k for up
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	newModel, _ = Update(newModel, msg)

	if newModel.MenuIndex != 0 {
		t.Errorf("MenuIndex = %d, want 0 after 'k'", newModel.MenuIndex)
	}
}

func TestDialogHeaderHeight(t *testing.T) {
	if DialogHeaderHeight <= 0 {
		t.Errorf("DialogHeaderHeight = %d, should be positive", DialogHeaderHeight)
	}
}

func TestQuitCommand(t *testing.T) {
	m := model.New("", 0)
	m.CtrlAPressed = true

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := Update(m, msg)

	if cmd == nil {
		t.Error("Quit command should return a tea.Cmd")
	}
}
