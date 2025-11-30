/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package serial

import (
	"errors"
	"time"

	"go.bug.st/serial"
)

// Common errors
var (
	ErrNoPortSelected = errors.New("no port selected")
	ErrPortNotOpen    = errors.New("port is not open")
)

const (
	// DefaultBaudRate is the standard baud rate for serial communication
	DefaultBaudRate = 9600

	// DefaultDataBits is the standard number of data bits
	DefaultDataBits = 8

	// DefaultStopBits is the standard number of stop bits
	DefaultStopBits = 1

	// DefaultParity is the default parity setting
	DefaultParity = "None"

	// DefaultLineEnding is the default line ending format
	DefaultLineEnding = "CRLF"

	// ReadTimeout is the timeout for serial port read operations
	ReadTimeout = 50 * time.Millisecond
)

// Config holds serial port configuration parameters.
type Config struct {
	Port       string // Serial port name (e.g., /dev/ttyUSB0, COM1)
	BaudRate   int    // Communication speed in bits per second
	DataBits   int    // Number of data bits (5-8)
	StopBits   int    // Number of stop bits (1-2)
	Parity     string // Parity mode: "None", "Even", "Odd", "Mark", "Space"
	LineEnding string // Line ending format: "NL", "CR", "CRLF", "None"
}

// DefaultConfig returns the default serial port configuration.
func DefaultConfig() Config {
	return Config{
		BaudRate:   DefaultBaudRate,
		DataBits:   DefaultDataBits,
		StopBits:   DefaultStopBits,
		Parity:     DefaultParity,
		LineEnding: DefaultLineEnding,
	}
}

// Connection manages the state and lifecycle of a serial port connection.
type Connection struct {
	Port      serial.Port // The underlying serial port handle
	Config    Config      // Port configuration parameters
	Connected bool        // Connection state flag
	LastError string      // Most recent error message
}

// NewConnection creates a new serial connection manager.
// If baudRate is 0 or negative, DefaultBaudRate will be used.
func NewConnection(port string, baudRate int) *Connection {
	cfg := DefaultConfig()
	cfg.Port = port

	if baudRate > 0 {
		cfg.BaudRate = baudRate
	}

	return &Connection{
		Config: cfg,
	}
}

// Connect opens the serial port with the configured parameters.
// Returns an error if connection fails. The error is also stored in LastError.
func (c *Connection) Connect() error {
	if c.Config.Port == "" {
		c.LastError = ErrNoPortSelected.Error()
		return ErrNoPortSelected
	}

	mode := &serial.Mode{
		BaudRate: c.Config.BaudRate,
		DataBits: c.Config.DataBits,
		StopBits: serial.StopBits(c.Config.StopBits),
		Parity:   c.parseParity(),
	}

	port, err := serial.Open(c.Config.Port, mode)
	if err != nil {
		c.LastError = err.Error()
		return err
	}

	// Set read timeout to prevent blocking indefinitely
	if err := port.SetReadTimeout(ReadTimeout); err != nil {
		_ = port.Close()
		c.LastError = err.Error()
		return err
	}

	c.Port = port
	c.Connected = true
	c.LastError = ""
	return nil
}

// parseParity converts the string parity setting to serial.Parity type.
func (c *Connection) parseParity() serial.Parity {
	switch c.Config.Parity {
	case "Even":
		return serial.EvenParity
	case "Odd":
		return serial.OddParity
	case "Mark":
		return serial.MarkParity
	case "Space":
		return serial.SpaceParity
	default:
		return serial.NoParity
	}
}

// Disconnect closes the serial port and resets connection state.
func (c *Connection) Disconnect() {
	if c.Port != nil {
		_ = c.Port.Close()
		c.Port = nil
	}
	c.Connected = false
}

// Write sends data to the serial port.
// Returns the number of bytes written and any error encountered.
func (c *Connection) Write(data []byte) (int, error) {
	if c.Port == nil {
		return 0, ErrPortNotOpen
	}
	return c.Port.Write(data)
}

// GetLineEndingBytes returns the line ending sequence based on configuration.
func (c *Connection) GetLineEndingBytes() string {
	switch c.Config.LineEnding {
	case "None":
		return ""
	case "NL":
		return "\n"
	case "CR":
		return "\r"
	case "CRLF":
		return "\r\n"
	default:
		return "\r\n"
	}
}

// GetLineEndingDisplay returns a human-readable description of the line ending.
func (c *Connection) GetLineEndingDisplay() string {
	switch c.Config.LineEnding {
	case "None":
		return "No Line Ending"
	case "NL":
		return "Newline (\\n)"
	case "CR":
		return "Carriage Return (\\r)"
	case "CRLF":
		return "Both NL & CR (\\r\\n)"
	default:
		return c.Config.LineEnding
	}
}

// BaudRates returns the list of commonly supported baud rates.
func BaudRates() []int {
	return []int{300, 1200, 2400, 4800, 9600, 19200, 38400, 57600, 115200, 230400}
}

// DataBitsOptions returns the valid data bits options.
func DataBitsOptions() []int {
	return []int{5, 6, 7, 8}
}

// StopBitsOptions returns the valid stop bits options.
func StopBitsOptions() []int {
	return []int{1, 2}
}

// LineEndingOptions returns the available line ending formats.
func LineEndingOptions() []string {
	return []string{"None", "NL", "CR", "CRLF"}
}
