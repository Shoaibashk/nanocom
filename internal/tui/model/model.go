/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package model

import (
	"io"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	bubblezone "github.com/lrstanley/bubblezone"
	"github.com/trzsz/trzsz-go/trzsz"

	"github.com/shoaibashk/nanocom/internal/tui/keys"
	"github.com/shoaibashk/nanocom/internal/tui/logging"
	"github.com/shoaibashk/nanocom/internal/tui/menu"
	"github.com/shoaibashk/nanocom/internal/tui/serial"
	"github.com/shoaibashk/nanocom/internal/tui/state"
	"github.com/shoaibashk/nanocom/internal/tui/transfer"
)

const (
	// StatusButtonZoneID is the zone ID for the clickable status bar button
	StatusButtonZoneID = "status-bar.button"

	// DefaultMaxLines is the default maximum number of terminal buffer lines
	DefaultMaxLines = 1000

	// DefaultInputCharLimit is the default character limit for input
	DefaultInputCharLimit = 256

	// DefaultInputWidth is the default width for the input field
	DefaultInputWidth = 80
)

// TransferSpinnerFrames defines the animation frames for the transfer spinner
var TransferSpinnerFrames = []string{"-", "\\", "|", "/"}

// Model represents the main TUI application state.
// It manages the terminal UI, serial connection, and all interactive components.
type Model struct {
	// View state management
	CurrentView state.ViewState

	// Terminal dimensions
	Width  int
	Height int

	// Serial connection management
	Serial *serial.Connection
	Reader *serial.Reader

	// Terminal buffer management
	TerminalBuffer []string
	MaxLines       int

	// User input field
	Input textinput.Model

	// Menu navigation state
	MenuIndex     int
	MenuItems     []menu.Item
	HelpMenuItems []menu.Item

	// Port selection state
	AvailablePorts []string
	PortIndex      int

	// Control key state
	CtrlAPressed bool

	// Key bindings configuration
	Keys keys.Map

	// File transfer state
	TransferState  transfer.State
	TransferBridge *transfer.Bridge

	// Logging state
	LoggingState logging.State

	// Transfer spinner state
	TransferPending          bool
	TransferPendingLabel     string
	TransferPendingDirection transfer.Direction
	TransferSpinnerIndex     int

	// UI interaction management
	Zone *bubblezone.Manager

	// Dialog drag state
	IsDragging   bool
	DragOffsetX  int
	DragOffsetY  int
	DialogX      int
	DialogY      int
	DialogWidth  int
	DialogHeight int
}

// New creates and initializes a new TUI model.
// If port is specified, the model will auto-connect on initialization.
// If baudRate is 0, the default baud rate (9600) will be used.
func New(port string, baudRate int) Model {
	ti := textinput.New()
	ti.Placeholder = "Type command here..."
	ti.Focus()
	ti.CharLimit = DefaultInputCharLimit
	ti.Width = DefaultInputWidth

	trzsz.SetAffectedByWindows(runtime.GOOS == "windows")

	m := Model{
		CurrentView:    state.ViewTerminal,
		Serial:         serial.NewConnection(port, baudRate),
		Reader:         serial.NewReader(),
		TerminalBuffer: make([]string, 0),
		MaxLines:       DefaultMaxLines,
		Input:          ti,
		MenuIndex:      0,
		PortIndex:      0,
		CtrlAPressed:   false,
		Keys:           keys.DefaultKeys(),
		TransferState:  transfer.NewState(),
		TransferBridge: transfer.NewBridge(),
		LoggingState:   logging.NewState(),
		Zone:           bubblezone.New(),
		MenuItems:      menu.DefaultMainMenuItems(),
		HelpMenuItems:  menu.DefaultHelpMenuItems(),
	}

	if m.Zone != nil {
		m.Zone.SetEnabled(true)
	}

	return m
}

// Init initializes the model and returns the initial commands.
// It starts the tick timer and refreshes available ports.
// If a port was specified via command line, it triggers auto-connect.
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		TickCmd(),
		m.RefreshPorts(),
	}

	// Auto-connect if port was specified via command line
	if m.Serial.Config.Port != "" {
		cmds = append(cmds, func() tea.Msg {
			return state.AutoConnectMsg{}
		})
	}

	return tea.Batch(cmds...)
}

// TickCmd returns a command that sends a tick message every second.
// Used for periodic UI updates like the transfer spinner animation.
func TickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return state.TickMsg(t)
	})
}

// RefreshPorts returns a command that refreshes the list of available serial ports.
func (m Model) RefreshPorts() tea.Cmd {
	return func() tea.Msg {
		ports, err := serial.GetPortsList()
		if err != nil {
			return nil
		}
		return state.PortsListMsg(ports)
	}
}

// ReadSerialCmd attempts a non-blocking read from the serial port channel.
// Returns immediately if no data is available.
func (m Model) ReadSerialCmd() tea.Cmd {
	return func() tea.Msg {
		if m.Reader.DataChan == nil {
			return nil
		}

		select {
		case data, ok := <-m.Reader.DataChan:
			if !ok {
				return nil
			}
			return state.SerialDataMsg(data)
		default:
			return nil
		}
	}
}

// WaitForSerialData blocks waiting for data from the serial channel.
// Used in the main read loop when connected.
func (m Model) WaitForSerialData() tea.Cmd {
	return func() tea.Msg {
		if m.Reader.DataChan == nil {
			return nil
		}

		data, ok := <-m.Reader.DataChan
		if !ok {
			return nil
		}
		return state.SerialDataMsg(data)
	}
}

// Connect establishes a connection to the configured serial port.
// On success, sets up the transfer bridge and starts the serial reader.
func (m *Model) Connect() {
	err := m.Serial.Connect()
	if err != nil {
		// Error is already stored in m.Serial.LastError
		return
	}

	m.AppendToBuffer("--- Connected to " + m.Serial.Config.Port + " ---")
	m.SetupTransferBridge()
	m.StartSerialReader()
}

// Disconnect closes the serial port connection and cleans up resources.
func (m *Model) Disconnect() {
	m.StopSerialReader()
	m.TeardownTransferBridge()
	m.Serial.Disconnect()
	m.AppendToBuffer("--- Disconnected ---")
}

// StartSerialReader starts the serial reader goroutine
func (m *Model) StartSerialReader() {
	reader := m.CurrentSerialReader()
	if reader != nil {
		m.Reader.Start(reader)
	}
}

// StopSerialReader stops the serial reader goroutine
func (m *Model) StopSerialReader() {
	m.Reader.Stop()
}

// CurrentSerialReader returns the current serial reader
func (m *Model) CurrentSerialReader() io.Reader {
	if m.TransferBridge.IsActive() {
		return m.TransferBridge.GetReader()
	}
	return m.Serial.Port
}

// WriteSerial writes data to the serial port.
// If a transfer bridge is active, writes through the bridge for protocol support.
// Returns early if data is empty or port is not connected.
func (m Model) WriteSerial(data []byte) {
	if len(data) == 0 {
		return
	}

	if m.TransferBridge.IsActive() {
		m.TransferBridge.Write(data)
		return
	}

	m.Serial.Write(data)
}

// SetupTransferBridge sets up the trzsz bridge
func (m *Model) SetupTransferBridge() {
	if m.Serial.Port == nil {
		return
	}
	m.TransferBridge.Setup(m.Serial.Port, m.CurrentTerminalColumns())
}

// TeardownTransferBridge tears down the trzsz bridge
func (m *Model) TeardownTransferBridge() {
	m.TransferBridge.Teardown()
	m.ClearTransferWait()
}

// CurrentTerminalColumns returns the current terminal columns
func (m Model) CurrentTerminalColumns() int32 {
	if m.Width <= 0 {
		return 80
	}
	return int32(m.Width)
}

// AppendToBuffer appends a line to the terminal buffer
func (m *Model) AppendToBuffer(line string) {
	m.TerminalBuffer = append(m.TerminalBuffer, line)
	if len(m.TerminalBuffer) > m.MaxLines {
		m.TerminalBuffer = m.TerminalBuffer[1:]
	}

	// Log the line if logging is enabled
	if m.LoggingState.IsEnabled() {
		_ = m.LoggingState.WriteLine(line)
	}
}

// BeginTransferWait starts the transfer spinner animation.
// Used to indicate a file transfer is waiting for remote action.
func (m *Model) BeginTransferWait(direction transfer.Direction, label string) {
	m.TransferPending = true
	m.TransferPendingDirection = direction
	m.TransferPendingLabel = label
	m.TransferSpinnerIndex = 0
}

// ClearTransferWait stops the transfer spinner animation.
func (m *Model) ClearTransferWait() {
	m.TransferPending = false
	m.TransferPendingLabel = ""
	m.TransferSpinnerIndex = 0
}

// AdvanceTransferSpinner moves to the next frame of the spinner animation.
func (m *Model) AdvanceTransferSpinner() {
	if !m.TransferPending || len(TransferSpinnerFrames) == 0 {
		return
	}
	m.TransferSpinnerIndex = (m.TransferSpinnerIndex + 1) % len(TransferSpinnerFrames)
}

// CurrentTransferStatus returns the current transfer status display string.
// Returns empty string if no transfer is pending.
func (m Model) CurrentTransferStatus() string {
	if !m.TransferPending || len(TransferSpinnerFrames) == 0 {
		return ""
	}

	frame := TransferSpinnerFrames[m.TransferSpinnerIndex%len(TransferSpinnerFrames)]
	label := m.TransferPendingLabel
	if label == "" {
		label = transfer.DirectionLabel(m.TransferPendingDirection)
	}

	return frame + " " + label
}

// GetLineEndingBytes returns the line ending bytes
func (m Model) GetLineEndingBytes() string {
	return m.Serial.GetLineEndingBytes()
}

// GetConnectLabel returns the connect/disconnect label
func (m Model) GetConnectLabel() string {
	if m.Serial.Connected {
		return "Disconnect"
	}
	return "Connect"
}
