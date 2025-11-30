/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package state

import (
	"testing"
	"time"
)

func TestViewStateConstants(t *testing.T) {
	// Verify view state constants are distinct
	views := []ViewState{
		ViewTerminal,
		ViewMainMenu,
		ViewSerialSetup,
		ViewHelp,
		ViewPortList,
		ViewSendFile,
		ViewReceiveFile,
	}

	seen := make(map[ViewState]bool)
	for _, v := range views {
		if seen[v] {
			t.Errorf("Duplicate ViewState value: %d", v)
		}
		seen[v] = true
	}
}

func TestViewStateOrdering(t *testing.T) {
	// ViewTerminal should be the first (default) state
	if ViewTerminal != 0 {
		t.Errorf("ViewTerminal = %d, want 0", ViewTerminal)
	}
}

func TestTickMsg(t *testing.T) {
	now := time.Now()
	msg := TickMsg(now)

	// TickMsg should be convertible to time.Time
	asTime := time.Time(msg)
	if !asTime.Equal(now) {
		t.Errorf("TickMsg time = %v, want %v", asTime, now)
	}
}

func TestSerialDataMsg(t *testing.T) {
	data := "Hello, World!"
	msg := SerialDataMsg(data)

	// SerialDataMsg should be convertible to string
	asString := string(msg)
	if asString != data {
		t.Errorf("SerialDataMsg = %q, want %q", asString, data)
	}
}

func TestSerialDataMsgEmpty(t *testing.T) {
	msg := SerialDataMsg("")

	if string(msg) != "" {
		t.Error("Empty SerialDataMsg should be empty string")
	}
}

func TestSerialErrorMsg(t *testing.T) {
	testErr := ErrTest{}
	msg := SerialErrorMsg{Err: testErr}

	if msg.Err == nil {
		t.Error("SerialErrorMsg.Err should not be nil")
	}
	if msg.Err.Error() != "test error" {
		t.Errorf("SerialErrorMsg.Err.Error() = %q, want 'test error'", msg.Err.Error())
	}
}

type ErrTest struct{}

func (e ErrTest) Error() string {
	return "test error"
}

func TestSerialErrorMsgNil(t *testing.T) {
	msg := SerialErrorMsg{Err: nil}

	if msg.Err != nil {
		t.Error("SerialErrorMsg.Err should be nil when created with nil")
	}
}

func TestPortsListMsg(t *testing.T) {
	ports := []string{"COM1", "COM2", "/dev/ttyUSB0"}
	msg := PortsListMsg(ports)

	if len(msg) != len(ports) {
		t.Errorf("PortsListMsg length = %d, want %d", len(msg), len(ports))
	}

	for i, port := range msg {
		if port != ports[i] {
			t.Errorf("PortsListMsg[%d] = %q, want %q", i, port, ports[i])
		}
	}
}

func TestPortsListMsgEmpty(t *testing.T) {
	msg := PortsListMsg([]string{})

	if len(msg) != 0 {
		t.Errorf("Empty PortsListMsg length = %d, want 0", len(msg))
	}
}

func TestPortsListMsgNil(t *testing.T) {
	msg := PortsListMsg(nil)

	if msg != nil {
		t.Error("Nil PortsListMsg should be nil")
	}
}

func TestAutoConnectMsg(t *testing.T) {
	msg := AutoConnectMsg{}

	// AutoConnectMsg is a marker type, just verify it can be created
	_ = msg
}

func TestViewStateNames(t *testing.T) {
	// Test that all view states are within expected range
	viewStates := []struct {
		name  string
		state ViewState
	}{
		{"ViewTerminal", ViewTerminal},
		{"ViewMainMenu", ViewMainMenu},
		{"ViewSerialSetup", ViewSerialSetup},
		{"ViewHelp", ViewHelp},
		{"ViewPortList", ViewPortList},
		{"ViewSendFile", ViewSendFile},
		{"ViewReceiveFile", ViewReceiveFile},
	}

	for _, vs := range viewStates {
		if vs.state < 0 {
			t.Errorf("%s should not be negative", vs.name)
		}
	}
}
