/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package transfer

import (
	"io"

	"github.com/trzsz/trzsz-go/trzsz"
	"go.bug.st/serial"
)

// SerialWriteCloser wraps a serial port for trzsz
type SerialWriteCloser struct {
	Port serial.Port
}

func (s SerialWriteCloser) Write(p []byte) (int, error) {
	if s.Port == nil {
		return 0, io.ErrClosedPipe
	}
	return s.Port.Write(p)
}

func (SerialWriteCloser) Close() error {
	return nil
}

// Bridge manages the trzsz filter for file transfers
type Bridge struct {
	Filter       *trzsz.TrzszFilter
	InputWriter  *io.PipeWriter
	OutputReader *io.PipeReader
}

// NewBridge creates a new trzsz bridge
func NewBridge() *Bridge {
	return &Bridge{}
}

// Setup initializes the trzsz bridge with the given serial port
func (b *Bridge) Setup(serialPort serial.Port, terminalColumns int32) {
	if b.Filter != nil || serialPort == nil {
		return
	}

	clientInReader, clientInWriter := io.Pipe()
	clientOutReader, clientOutWriter := io.Pipe()

	b.Filter = trzsz.NewTrzszFilter(
		clientInReader,
		clientOutWriter,
		SerialWriteCloser{Port: serialPort},
		serialPort,
		trzsz.TrzszOptions{
			TerminalColumns: terminalColumns,
			EnableZmodem:    true,
		},
	)
	b.InputWriter = clientInWriter
	b.OutputReader = clientOutReader
}

// Teardown cleans up the trzsz bridge
func (b *Bridge) Teardown() {
	if b.Filter != nil {
		b.Filter.StopTransferringFiles(true)
		b.Filter.ResetTerminal()
		b.Filter = nil
	}
	if b.InputWriter != nil {
		_ = b.InputWriter.Close()
		b.InputWriter = nil
	}
	if b.OutputReader != nil {
		_ = b.OutputReader.Close()
		b.OutputReader = nil
	}
}

// IsActive returns whether the bridge is active
func (b *Bridge) IsActive() bool {
	return b.Filter != nil
}

// GetReader returns the output reader for receiving data
func (b *Bridge) GetReader() io.Reader {
	return b.OutputReader
}

// Write writes data through the bridge
func (b *Bridge) Write(data []byte) {
	if len(data) == 0 {
		return
	}
	if b.InputWriter != nil {
		_, _ = b.InputWriter.Write(data)
	}
}

// SetTerminalColumns updates the terminal columns
func (b *Bridge) SetTerminalColumns(columns int32) {
	if b.Filter != nil {
		b.Filter.SetTerminalColumns(columns)
	}
}

// IsTransferring returns whether a file transfer is in progress
func (b *Bridge) IsTransferring() bool {
	return b.Filter != nil && b.Filter.IsTransferringFiles()
}

// SetDefaultUploadPath sets the default upload path
func (b *Bridge) SetDefaultUploadPath(path string) {
	if b.Filter != nil {
		b.Filter.SetDefaultUploadPath(path)
	}
}

// SetDefaultDownloadPath sets the default download path
func (b *Bridge) SetDefaultDownloadPath(path string) {
	if b.Filter != nil {
		b.Filter.SetDefaultDownloadPath(path)
	}
}

// OneTimeUpload initiates a one-time upload
func (b *Bridge) OneTimeUpload(files []string) (<-chan error, error) {
	if b.Filter == nil {
		return nil, io.ErrClosedPipe
	}
	return b.Filter.OneTimeUpload(files)
}
