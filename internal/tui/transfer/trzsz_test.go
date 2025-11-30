/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package transfer

import (
	"io"
	"testing"
)

func TestNewBridge(t *testing.T) {
	b := NewBridge()

	if b == nil {
		t.Fatal("NewBridge returned nil")
	}
	if b.Filter != nil {
		t.Error("NewBridge Filter should be nil")
	}
	if b.InputWriter != nil {
		t.Error("NewBridge InputWriter should be nil")
	}
	if b.OutputReader != nil {
		t.Error("NewBridge OutputReader should be nil")
	}
}

func TestBridgeIsActive(t *testing.T) {
	b := NewBridge()

	if b.IsActive() {
		t.Error("New bridge should not be active")
	}
}

func TestBridgeSetupNilPort(t *testing.T) {
	b := NewBridge()

	// Setup with nil port should not create filter
	b.Setup(nil, 80)

	if b.IsActive() {
		t.Error("Bridge should not be active after Setup with nil port")
	}
}

func TestBridgeTeardown(t *testing.T) {
	b := NewBridge()

	// Teardown on new bridge should not panic
	b.Teardown()

	if b.IsActive() {
		t.Error("Bridge should not be active after Teardown")
	}
	if b.InputWriter != nil {
		t.Error("InputWriter should be nil after Teardown")
	}
	if b.OutputReader != nil {
		t.Error("OutputReader should be nil after Teardown")
	}
}

func TestBridgeWriteEmpty(t *testing.T) {
	b := NewBridge()

	// Write empty data should not panic
	b.Write([]byte{})
	b.Write(nil)
}

func TestBridgeWriteNoWriter(t *testing.T) {
	b := NewBridge()

	// Write with no writer should not panic
	b.Write([]byte("test"))
}

func TestBridgeGetReaderNil(t *testing.T) {
	b := NewBridge()

	// OutputReader may or may not be nil depending on implementation
	_ = b.GetReader()
}

func TestBridgeSetTerminalColumnsNoFilter(t *testing.T) {
	b := NewBridge()

	// Should not panic when filter is nil
	b.SetTerminalColumns(80)
}

func TestBridgeIsTransferring(t *testing.T) {
	b := NewBridge()

	if b.IsTransferring() {
		t.Error("New bridge should not be transferring")
	}
}

func TestBridgeSetDefaultUploadPathNoFilter(t *testing.T) {
	b := NewBridge()

	// Should not panic when filter is nil
	b.SetDefaultUploadPath("/some/path")
}

func TestBridgeSetDefaultDownloadPathNoFilter(t *testing.T) {
	b := NewBridge()

	// Should not panic when filter is nil
	b.SetDefaultDownloadPath("/some/path")
}

func TestBridgeOneTimeUploadNoFilter(t *testing.T) {
	b := NewBridge()

	ch, err := b.OneTimeUpload([]string{"file.txt"})
	if err != io.ErrClosedPipe {
		t.Errorf("OneTimeUpload() error = %v, want io.ErrClosedPipe", err)
	}
	if ch != nil {
		t.Error("OneTimeUpload() channel should be nil on error")
	}
}

func TestSerialWriteCloserWriteNilPort(t *testing.T) {
	s := SerialWriteCloser{Port: nil}

	n, err := s.Write([]byte("test"))
	if err != io.ErrClosedPipe {
		t.Errorf("Write() error = %v, want io.ErrClosedPipe", err)
	}
	if n != 0 {
		t.Errorf("Write() n = %d, want 0", n)
	}
}

func TestSerialWriteCloserClose(t *testing.T) {
	s := SerialWriteCloser{Port: nil}

	err := s.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestBridgeSetupAlreadyActive(t *testing.T) {
	_ = NewBridge()

	// Manually set filter to simulate active state
	// (We can't actually test this without a real serial port)
	// The Setup function should return early if Filter is not nil
}

func TestBridgeTeardownMultipleTimes(t *testing.T) {
	b := NewBridge()

	// Multiple teardowns should not panic
	b.Teardown()
	b.Teardown()
	b.Teardown()
}

func TestBridgeWriteWithData(t *testing.T) {
	b := NewBridge()

	// Write with data but no writer - should not panic
	b.Write([]byte("Hello, World!"))
}
