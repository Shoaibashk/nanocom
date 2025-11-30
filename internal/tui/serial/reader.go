/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package serial

import (
	"io"
	"time"
)

const (
	// DataChannelBufferSize is the buffer size for the data channel
	DataChannelBufferSize = 100

	// ReadBufferSize is the size of the read buffer for serial data
	ReadBufferSize = 2048

	// ReadRetryDelay is the delay before retrying after EOF
	ReadRetryDelay = 50 * time.Millisecond
)

// Reader manages asynchronous reading from a serial port.
// It runs a background goroutine that continuously reads data and
// sends it to a channel for consumption by the main UI loop.
type Reader struct {
	DataChan     chan string   // Channel for sending received data
	StopChan     chan struct{} // Channel for signaling goroutine to stop
	Reading      bool          // Flag indicating if reader is active
	SerialBuffer string        // Buffer for incomplete lines (no newline yet)
}

// NewReader creates a new serial reader in an initialized but inactive state.
func NewReader() *Reader {
	return &Reader{}
}

// Start begins asynchronously reading from the provided reader.
// It spawns a goroutine that continuously reads data and sends it to DataChan.
// Does nothing if already reading or if reader is nil.
func (r *Reader) Start(reader io.Reader) {
	if r.Reading || reader == nil {
		return
	}

	r.DataChan = make(chan string, DataChannelBufferSize)
	r.StopChan = make(chan struct{})
	r.Reading = true

	go r.readLoop(reader)
}

// readLoop is the main read goroutine that continuously reads from the serial port.
func (r *Reader) readLoop(rd io.Reader) {
	buf := make([]byte, ReadBufferSize)

	for {
		// Check if we should stop
		select {
		case <-r.StopChan:
			return
		default:
		}

		// Attempt to read data
		n, err := rd.Read(buf)
		if n > 0 {
			payload := string(buf[:n])

			// Try to send data, but also check for stop signal
			select {
			case r.DataChan <- payload:
			case <-r.StopChan:
				return
			}
		}

		// Handle read errors (EOF is common and expected)
		if err != nil {
			if err == io.EOF {
				time.Sleep(ReadRetryDelay)
			}
			continue
		}
	}
}

// Stop halts the serial reader goroutine and cleans up resources.
// It closes channels, resets state, and clears any buffered data.
// Safe to call multiple times.
func (r *Reader) Stop() {
	if r.StopChan != nil {
		close(r.StopChan)
		r.StopChan = nil
	}

	if r.DataChan != nil {
		close(r.DataChan)
		r.DataChan = nil
	}

	r.Reading = false
	r.SerialBuffer = ""
}

// ClearBuffer clears any incomplete buffered data.
// Used when resetting or reconnecting.
func (r *Reader) ClearBuffer() {
	r.SerialBuffer = ""
}
