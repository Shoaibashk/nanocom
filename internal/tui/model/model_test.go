/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package model

import (
	"testing"

	"github.com/shoaibashk/nanocom/internal/tui/state"
	"github.com/shoaibashk/nanocom/internal/tui/transfer"
)

func TestNew(t *testing.T) {
	m := New("", 0)

	if m.CurrentView != state.ViewTerminal {
		t.Errorf("New().CurrentView = %v, want ViewTerminal", m.CurrentView)
	}
	if m.Serial == nil {
		t.Error("New().Serial should not be nil")
	}
	if m.Reader == nil {
		t.Error("New().Reader should not be nil")
	}
	if m.MaxLines != DefaultMaxLines {
		t.Errorf("New().MaxLines = %d, want %d", m.MaxLines, DefaultMaxLines)
	}
	if m.TransferBridge == nil {
		t.Error("New().TransferBridge should not be nil")
	}
	if m.Zone == nil {
		t.Error("New().Zone should not be nil")
	}
	if len(m.MenuItems) == 0 {
		t.Error("New().MenuItems should not be empty")
	}
	if len(m.HelpMenuItems) == 0 {
		t.Error("New().HelpMenuItems should not be empty")
	}
}

func TestNewWithPort(t *testing.T) {
	m := New("COM1", 115200)

	if m.Serial.Config.Port != "COM1" {
		t.Errorf("Config.Port = %q, want 'COM1'", m.Serial.Config.Port)
	}
	if m.Serial.Config.BaudRate != 115200 {
		t.Errorf("Config.BaudRate = %d, want 115200", m.Serial.Config.BaudRate)
	}
}

func TestAppendToBuffer(t *testing.T) {
	m := New("", 0)

	m.AppendToBuffer("Line 1")
	m.AppendToBuffer("Line 2")

	if len(m.TerminalBuffer) != 2 {
		t.Errorf("TerminalBuffer length = %d, want 2", len(m.TerminalBuffer))
	}
	if m.TerminalBuffer[0] != "Line 1" {
		t.Errorf("TerminalBuffer[0] = %q, want 'Line 1'", m.TerminalBuffer[0])
	}
	if m.TerminalBuffer[1] != "Line 2" {
		t.Errorf("TerminalBuffer[1] = %q, want 'Line 2'", m.TerminalBuffer[1])
	}
}

func TestAppendToBufferMaxLines(t *testing.T) {
	m := New("", 0)
	m.MaxLines = 5

	for i := 0; i < 10; i++ {
		m.AppendToBuffer("Line")
	}

	if len(m.TerminalBuffer) != 5 {
		t.Errorf("TerminalBuffer length = %d, want 5 (maxed)", len(m.TerminalBuffer))
	}
}

func TestCurrentTerminalColumns(t *testing.T) {
	m := New("", 0)

	// Default should return 80
	if cols := m.CurrentTerminalColumns(); cols != 80 {
		t.Errorf("CurrentTerminalColumns() = %d, want 80 (default)", cols)
	}

	m.Width = 120
	if cols := m.CurrentTerminalColumns(); cols != 120 {
		t.Errorf("CurrentTerminalColumns() = %d, want 120", cols)
	}

	m.Width = -1
	if cols := m.CurrentTerminalColumns(); cols != 80 {
		t.Errorf("CurrentTerminalColumns() with negative width = %d, want 80", cols)
	}
}

func TestGetConnectLabel(t *testing.T) {
	m := New("", 0)

	if label := m.GetConnectLabel(); label != "Connect" {
		t.Errorf("GetConnectLabel() = %q, want 'Connect'", label)
	}

	m.Serial.Connected = true
	if label := m.GetConnectLabel(); label != "Disconnect" {
		t.Errorf("GetConnectLabel() = %q, want 'Disconnect'", label)
	}
}

func TestGetLineEndingBytes(t *testing.T) {
	m := New("", 0)

	ending := m.GetLineEndingBytes()
	if ending != "\r\n" {
		t.Errorf("GetLineEndingBytes() = %q, want '\\r\\n' (default)", ending)
	}
}

func TestBeginTransferWait(t *testing.T) {
	m := New("", 0)

	m.BeginTransferWait(transfer.DirectionSend, "Sending file.txt")

	if !m.TransferPending {
		t.Error("TransferPending should be true")
	}
	if m.TransferPendingDirection != transfer.DirectionSend {
		t.Error("TransferPendingDirection should be DirectionSend")
	}
	if m.TransferPendingLabel != "Sending file.txt" {
		t.Errorf("TransferPendingLabel = %q, want 'Sending file.txt'", m.TransferPendingLabel)
	}
	if m.TransferSpinnerIndex != 0 {
		t.Errorf("TransferSpinnerIndex = %d, want 0", m.TransferSpinnerIndex)
	}
}

func TestClearTransferWait(t *testing.T) {
	m := New("", 0)
	m.TransferPending = true
	m.TransferPendingLabel = "test"
	m.TransferSpinnerIndex = 5

	m.ClearTransferWait()

	if m.TransferPending {
		t.Error("TransferPending should be false")
	}
	if m.TransferPendingLabel != "" {
		t.Error("TransferPendingLabel should be empty")
	}
	if m.TransferSpinnerIndex != 0 {
		t.Error("TransferSpinnerIndex should be 0")
	}
}

func TestAdvanceTransferSpinner(t *testing.T) {
	m := New("", 0)
	m.TransferPending = true

	for i := 0; i < len(TransferSpinnerFrames)*2; i++ {
		m.AdvanceTransferSpinner()
	}

	// Should wrap around
	if m.TransferSpinnerIndex >= len(TransferSpinnerFrames) {
		t.Errorf("TransferSpinnerIndex = %d, should wrap around", m.TransferSpinnerIndex)
	}
}

func TestAdvanceTransferSpinnerNotPending(t *testing.T) {
	m := New("", 0)
	m.TransferPending = false
	m.TransferSpinnerIndex = 0

	m.AdvanceTransferSpinner()

	if m.TransferSpinnerIndex != 0 {
		t.Error("TransferSpinnerIndex should not change when not pending")
	}
}

func TestCurrentTransferStatus(t *testing.T) {
	m := New("", 0)

	// Not pending
	if status := m.CurrentTransferStatus(); status != "" {
		t.Errorf("CurrentTransferStatus() = %q, want empty when not pending", status)
	}

	// Pending
	m.TransferPending = true
	m.TransferPendingLabel = "Waiting..."

	status := m.CurrentTransferStatus()
	if status == "" {
		t.Error("CurrentTransferStatus() should not be empty when pending")
	}
}

func TestTransferSpinnerFrames(t *testing.T) {
	if len(TransferSpinnerFrames) == 0 {
		t.Fatal("TransferSpinnerFrames should not be empty")
	}

	for i, frame := range TransferSpinnerFrames {
		if frame == "" {
			t.Errorf("TransferSpinnerFrames[%d] should not be empty", i)
		}
	}
}

func TestConstants(t *testing.T) {
	if StatusButtonZoneID == "" {
		t.Error("StatusButtonZoneID should not be empty")
	}
	if DefaultMaxLines <= 0 {
		t.Errorf("DefaultMaxLines = %d, should be positive", DefaultMaxLines)
	}
	if DefaultInputCharLimit <= 0 {
		t.Errorf("DefaultInputCharLimit = %d, should be positive", DefaultInputCharLimit)
	}
	if DefaultInputWidth <= 0 {
		t.Errorf("DefaultInputWidth = %d, should be positive", DefaultInputWidth)
	}
}

func TestModelInit(t *testing.T) {
	m := New("", 0)

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestModelInitWithPort(t *testing.T) {
	m := New("COM1", 9600)

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() with port should return a command")
	}
}

func TestRefreshPorts(t *testing.T) {
	m := New("", 0)

	cmd := m.RefreshPorts()
	if cmd == nil {
		t.Error("RefreshPorts() should return a command")
	}

	// Execute the command
	msg := cmd()
	// msg could be nil or PortsListMsg
	_ = msg
}

func TestWriteSerialEmpty(t *testing.T) {
	m := New("", 0)

	// Should not panic
	m.WriteSerial([]byte{})
	m.WriteSerial(nil)
}

func TestWriteSerialNotConnected(t *testing.T) {
	m := New("", 0)

	// Should not panic when not connected
	m.WriteSerial([]byte("test"))
}

func TestTickCmd(t *testing.T) {
	cmd := TickCmd()
	if cmd == nil {
		t.Error("TickCmd() should return a command")
	}
}

func TestReadSerialCmdNoChannel(t *testing.T) {
	m := New("", 0)
	m.Reader.DataChan = nil

	cmd := m.ReadSerialCmd()
	if cmd == nil {
		t.Error("ReadSerialCmd() should return a command")
	}

	msg := cmd()
	if msg != nil {
		t.Error("ReadSerialCmd() with nil channel should return nil msg")
	}
}

func TestWaitForSerialDataNoChannel(t *testing.T) {
	m := New("", 0)
	m.Reader.DataChan = nil

	cmd := m.WaitForSerialData()
	if cmd == nil {
		t.Error("WaitForSerialData() should return a command")
	}

	msg := cmd()
	if msg != nil {
		t.Error("WaitForSerialData() with nil channel should return nil msg")
	}
}
