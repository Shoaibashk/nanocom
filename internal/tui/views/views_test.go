/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package views

import (
	"strings"
	"testing"

	"github.com/shoaibashk/nanocom/internal/tui/model"
	"github.com/shoaibashk/nanocom/internal/tui/state"
)

func TestViewTerminal(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string")
	}
}

func TestViewMainMenu(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewMainMenu

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for main menu")
	}
	if !strings.Contains(view, "nanocom") {
		t.Error("Main menu should contain 'nanocom'")
	}
}

func TestViewSerialSetup(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewSerialSetup

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for serial setup")
	}
	if !strings.Contains(view, "Serial") {
		t.Error("Serial setup should contain 'Serial'")
	}
}

func TestViewHelp(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewHelp

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for help")
	}
	if !strings.Contains(view, "Help") {
		t.Error("Help view should contain 'Help'")
	}
}

func TestViewPortList(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewPortList
	m.AvailablePorts = []string{"COM1", "COM2"}

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for port list")
	}
	if !strings.Contains(view, "Serial Port") {
		t.Error("Port list should contain 'Serial Port'")
	}
}

func TestViewPortListEmpty(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewPortList
	m.AvailablePorts = []string{}

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for empty port list")
	}
	if !strings.Contains(view, "No serial ports") {
		t.Error("Empty port list should indicate no ports found")
	}
}

func TestViewSendFile(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewSendFile

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for send file")
	}
	if !strings.Contains(view, "Send") {
		t.Error("Send file view should contain 'Send'")
	}
}

func TestViewReceiveFile(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.CurrentView = state.ViewReceiveFile

	view := View(m)

	if view == "" {
		t.Error("View() should return non-empty string for receive file")
	}
	if !strings.Contains(view, "Receive") {
		t.Error("Receive file view should contain 'Receive'")
	}
}

func TestViewNotInitialized(t *testing.T) {
	m := model.New("", 0)
	// Width and Height are 0

	view := View(m)

	if !strings.Contains(view, "Initializing") {
		t.Error("View with 0 dimensions should show 'Initializing'")
	}
}

func TestRenderTerminalView(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	view := RenderTerminalView(m)

	if view == "" {
		t.Error("RenderTerminalView should return non-empty string")
	}
}

func TestRenderStatusBar(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	status := RenderStatusBar(m)

	if status == "" {
		t.Error("RenderStatusBar should return non-empty string")
	}
}

func TestRenderStatusBarConnected(t *testing.T) {
	m := model.New("COM1", 9600)
	m.Width = 80
	m.Height = 24
	m.Serial.Connected = true

	status := RenderStatusBar(m)

	if !strings.Contains(status, "ONLINE") && !strings.Contains(status, "Disconnect") {
		t.Error("Status bar for connected state should show ONLINE or Disconnect")
	}
}

func TestRenderStatusBarDisconnected(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.Serial.Connected = false

	status := RenderStatusBar(m)

	if !strings.Contains(status, "OFFLINE") && !strings.Contains(status, "Connect") {
		t.Error("Status bar for disconnected state should show OFFLINE or Connect")
	}
}

func TestRenderStatusBarWithLogging(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	tmpDir := t.TempDir()
	logPath := tmpDir + "/test.log"
	err := m.LoggingState.Start(logPath)
	if err != nil {
		t.Fatalf("Failed to start logging: %v", err)
	}
	defer m.LoggingState.Stop()

	status := RenderStatusBar(m)

	// The status bar should render correctly when logging is enabled
	if status == "" {
		t.Error("Status bar should not be empty when logging is enabled")
	}
}

func TestRenderStatusBarWithTransfer(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.TransferPending = true
	m.TransferPendingLabel = "Sending..."

	status := RenderStatusBar(m)

	if status == "" {
		t.Error("Status bar should render with transfer pending")
	}
}

func TestRenderMainMenu(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	menu := RenderMainMenu(m)

	if menu == "" {
		t.Error("RenderMainMenu should return non-empty string")
	}
	if !strings.Contains(menu, "nanocom") {
		t.Error("Main menu should contain 'nanocom'")
	}
}

func TestRenderSerialSetup(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	setup := RenderSerialSetup(m)

	if setup == "" {
		t.Error("RenderSerialSetup should return non-empty string")
	}
	if !strings.Contains(setup, "Serial") {
		t.Error("Serial setup should contain 'Serial'")
	}
}

func TestRenderSerialSetupWithError(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.Serial.LastError = "Connection failed"

	setup := RenderSerialSetup(m)

	if !strings.Contains(setup, "Connection failed") {
		t.Error("Serial setup should show error message")
	}
}

func TestRenderHelpMenu(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	help := RenderHelpMenu(m)

	if help == "" {
		t.Error("RenderHelpMenu should return non-empty string")
	}
	if !strings.Contains(help, "Help") {
		t.Error("Help menu should contain 'Help'")
	}
}

func TestRenderPortList(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.AvailablePorts = []string{"COM1", "COM2", "COM3"}
	m.PortIndex = 1

	list := RenderPortList(m)

	if list == "" {
		t.Error("RenderPortList should return non-empty string")
	}
	if !strings.Contains(list, "COM1") || !strings.Contains(list, "COM2") {
		t.Error("Port list should contain available ports")
	}
}

func TestRenderSendFileView(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	view := RenderSendFileView(m)

	if view == "" {
		t.Error("RenderSendFileView should return non-empty string")
	}
	if !strings.Contains(view, "Send") {
		t.Error("Send file view should contain 'Send'")
	}
}

func TestRenderReceiveFileView(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	view := RenderReceiveFileView(m)

	if view == "" {
		t.Error("RenderReceiveFileView should return non-empty string")
	}
	if !strings.Contains(view, "Receive") {
		t.Error("Receive file view should contain 'Receive'")
	}
}

func TestGetEstimatedDialogSize(t *testing.T) {
	tests := []struct {
		view state.ViewState
	}{
		{state.ViewMainMenu},
		{state.ViewSerialSetup},
		{state.ViewHelp},
		{state.ViewPortList},
		{state.ViewSendFile},
		{state.ViewReceiveFile},
		{state.ViewTerminal},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			m := model.New("", 0)
			m.CurrentView = tt.view
			m.AvailablePorts = []string{"COM1"}

			width, height := GetEstimatedDialogSize(m)

			if width <= 0 {
				t.Errorf("Dialog width = %d, should be positive", width)
			}
			if height <= 0 {
				t.Errorf("Dialog height = %d, should be positive", height)
			}
		})
	}
}

func TestGetDialogX(t *testing.T) {
	m := model.New("", 0)
	m.Width = 100
	m.Height = 50

	x := GetDialogX(m)

	// Should be centered
	if x < 0 || x >= m.Width {
		t.Errorf("DialogX = %d, should be within bounds", x)
	}

	// With explicit position
	m.DialogX = 10
	x = GetDialogX(m)
	if x != 10 {
		t.Errorf("DialogX = %d, want 10 (explicit)", x)
	}
}

func TestGetDialogY(t *testing.T) {
	m := model.New("", 0)
	m.Width = 100
	m.Height = 50

	y := GetDialogY(m)

	// Should be centered
	if y < 0 || y >= m.Height {
		t.Errorf("DialogY = %d, should be within bounds", y)
	}

	// With explicit position
	m.DialogY = 5
	y = GetDialogY(m)
	if y != 5 {
		t.Errorf("DialogY = %d, want 5 (explicit)", y)
	}
}

func TestOverlayAtWithSize(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24

	bg := strings.Repeat(strings.Repeat(" ", 80)+"\n", 24)
	overlay := "OVERLAY"

	result := OverlayAtWithSize(m, bg, overlay, 10, 10, 7, 1)

	if result == "" {
		t.Error("OverlayAtWithSize should return non-empty string")
	}
	if !strings.Contains(result, "OVERLAY") {
		t.Error("Result should contain overlay text")
	}
}

func TestTerminalBufferContent(t *testing.T) {
	m := model.New("", 0)
	m.Width = 80
	m.Height = 24
	m.AppendToBuffer("Test line 1")
	m.AppendToBuffer("Test line 2")
	m.AppendToBuffer("> User input")
	m.AppendToBuffer("--- System message ---")

	view := RenderTerminalView(m)

	if !strings.Contains(view, "Test line 1") {
		t.Error("Terminal view should contain buffer content")
	}
}
