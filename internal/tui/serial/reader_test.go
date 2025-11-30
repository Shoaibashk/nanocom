/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package serial

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	r := NewReader()

	if r == nil {
		t.Fatal("NewReader returned nil")
	}
	if r.Reading {
		t.Error("NewReader should not be reading")
	}
	if r.DataChan != nil {
		t.Error("NewReader DataChan should be nil")
	}
	if r.StopChan != nil {
		t.Error("NewReader StopChan should be nil")
	}
	if r.SerialBuffer != "" {
		t.Error("NewReader SerialBuffer should be empty")
	}
}

func TestStartStop(t *testing.T) {
	r := NewReader()

	// Create a simple reader
	testData := "Hello, World!"
	reader := strings.NewReader(testData)

	r.Start(reader)

	if !r.Reading {
		t.Error("Reader should be reading after Start")
	}
	if r.DataChan == nil {
		t.Error("DataChan should not be nil after Start")
	}
	if r.StopChan == nil {
		t.Error("StopChan should not be nil after Start")
	}

	// Give time for the goroutine to read
	time.Sleep(100 * time.Millisecond)

	r.Stop()

	if r.Reading {
		t.Error("Reader should not be reading after Stop")
	}
	if r.DataChan != nil {
		t.Error("DataChan should be nil after Stop")
	}
	if r.StopChan != nil {
		t.Error("StopChan should be nil after Stop")
	}
}

func TestStartWithNilReader(t *testing.T) {
	r := NewReader()

	r.Start(nil)

	if r.Reading {
		t.Error("Reader should not start with nil reader")
	}
}

func TestStartAlreadyReading(t *testing.T) {
	r := NewReader()
	reader := strings.NewReader("test")

	r.Start(reader)
	defer r.Stop()

	// Try to start again
	originalChan := r.DataChan
	r.Start(reader)

	// Channel should be the same (not restarted)
	if r.DataChan != originalChan {
		t.Error("DataChan should not change when starting already-reading reader")
	}
}

func TestStopMultipleTimes(t *testing.T) {
	r := NewReader()
	reader := strings.NewReader("test")

	r.Start(reader)
	r.Stop()
	r.Stop() // Should not panic
	r.Stop() // Should not panic
}

func TestStopNotStarted(t *testing.T) {
	r := NewReader()

	// Should not panic
	r.Stop()
}

func TestClearBuffer(t *testing.T) {
	r := NewReader()
	r.SerialBuffer = "some buffered data"

	r.ClearBuffer()

	if r.SerialBuffer != "" {
		t.Error("SerialBuffer should be empty after ClearBuffer")
	}
}

func TestReadData(t *testing.T) {
	r := NewReader()
	testData := "Hello, World!"
	reader := strings.NewReader(testData)

	r.Start(reader)
	defer r.Stop()

	// Wait for data
	select {
	case data := <-r.DataChan:
		if data != testData {
			t.Errorf("Received data = %q, want %q", data, testData)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for data")
	}
}

type slowReader struct {
	data  string
	sent  bool
	delay time.Duration
}

func (r *slowReader) Read(p []byte) (int, error) {
	if r.sent {
		time.Sleep(r.delay)
		return 0, io.EOF
	}
	r.sent = true
	n := copy(p, r.data)
	return n, nil
}

func TestReadWithDelay(t *testing.T) {
	r := NewReader()
	reader := &slowReader{
		data:  "delayed data",
		delay: 100 * time.Millisecond,
	}

	r.Start(reader)
	defer r.Stop()

	select {
	case data := <-r.DataChan:
		if data != "delayed data" {
			t.Errorf("Received data = %q, want 'delayed data'", data)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for data")
	}
}

func TestReaderConstants(t *testing.T) {
	if DataChannelBufferSize <= 0 {
		t.Errorf("DataChannelBufferSize = %d, should be positive", DataChannelBufferSize)
	}
	if ReadBufferSize <= 0 {
		t.Errorf("ReadBufferSize = %d, should be positive", ReadBufferSize)
	}
	if ReadRetryDelay <= 0 {
		t.Errorf("ReadRetryDelay = %v, should be positive", ReadRetryDelay)
	}
}

type infiniteReader struct {
	count int
}

func (r *infiniteReader) Read(p []byte) (int, error) {
	if r.count >= 10 {
		return 0, io.EOF
	}
	r.count++
	copy(p, "x")
	return 1, nil
}

func TestMultipleReads(t *testing.T) {
	r := NewReader()
	reader := &infiniteReader{}

	r.Start(reader)

	// Read several chunks
	received := 0
	timeout := time.After(2 * time.Second)

	for received < 5 {
		select {
		case <-r.DataChan:
			received++
		case <-timeout:
			r.Stop()
			t.Fatalf("Timeout: only received %d chunks", received)
			return
		}
	}

	r.Stop()

	if received < 5 {
		t.Errorf("Received %d chunks, want at least 5", received)
	}
}
