/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	"io"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	bubblezone "github.com/lrstanley/bubblezone"
	"github.com/trzsz/trzsz-go/trzsz"
	"go.bug.st/serial"
)

// ViewState represents the current view/mode of the TUI
type ViewState int

const (
	ViewTerminal ViewState = iota
	ViewMainMenu
	ViewSerialSetup
	ViewHelp
	ViewPortList
	ViewSendFile
	ViewReceiveFile
)

const statusButtonZoneID = "status-bar.button"

// tickMsg for periodic updates (like time)
type tickMsg time.Time

// serialDataMsg for incoming serial data
type serialDataMsg string

var transferSpinnerFrames = []string{"-", "\\", "|", "/"}

// Model represents the main TUI model
type Model struct {
	// View state
	currentView ViewState

	// Terminal dimensions
	width  int
	height int

	// Serial settings
	port       string
	baudRate   int
	dataBits   int
	stopBits   int
	parity     string
	lineEnding string // "NL", "CR", "CRLF", "None"

	// Serial port handle
	serialPort serial.Port

	// Terminal buffer
	terminalBuffer []string
	maxLines       int

	// Input for sending data
	input textinput.Model

	// Menu state
	menuIndex     int
	menuItems     []menuItem
	helpMenuItems []menuItem

	// Available ports
	availablePorts []string
	portIndex      int

	// Status
	connected bool
	lastError string

	// Control key pressed
	ctrlAPressed bool

	// Buffer for incomplete serial data (no newline yet)
	serialBuffer string

	// Serial reader control
	serialChan   chan string
	stopReadChan chan struct{}
	reading      bool

	// Key bindings
	keys keyMap

	// File transfer state
	transferState TransferState

	// Logging state
	loggingState LoggingState

	transferPending          bool
	transferPendingLabel     string
	transferPendingDirection TransferDirection
	transferSpinnerIndex     int

	// trzsz filter bridge for file transfers
	trzFilter          *trzsz.TrzszFilter
	filterInputWriter  *io.PipeWriter
	filterOutputReader *io.PipeReader

	// Bubble Zone manager for clickable regions
	zone *bubblezone.Manager

	// Mouse dragging state for dialog windows
	isDragging   bool
	dragOffsetX  int
	dragOffsetY  int
	dialogX      int
	dialogY      int
	dialogWidth  int
	dialogHeight int
}

// Serial setup menu item indices
const (
	serialSetupPort = iota
	serialSetupBaudRate
	serialSetupDataBits
	serialSetupStopBits
	serialSetupLineEnding
	serialSetupConnect
	serialSetupMenuItemCount
)

// menuItem represents a menu option
type menuItem struct {
	key   string
	label string
}

// keyMap defines the key bindings
type keyMap struct {
	Quit       key.Binding
	Menu       key.Binding
	Enter      key.Binding
	Up         key.Binding
	Down       key.Binding
	Escape     key.Binding
	CtrlA      key.Binding
	Help       key.Binding
	Connect    key.Binding
	Disconnect key.Binding
}

var defaultKeys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Menu: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("^A", "menu"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	CtrlA: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("^A", "command"),
	),
	Help: key.NewBinding(
		key.WithKeys("z"),
		key.WithHelp("Z", "help"),
	),
	Connect: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("C", "connect"),
	),
	Disconnect: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("D", "disconnect"),
	),
}

// NewModel creates a new TUI model
func NewModel(port string, baudRate int) Model {
	ti := textinput.New()
	ti.Placeholder = "Type command here..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80

	trzsz.SetAffectedByWindows(runtime.GOOS == "windows")

	// Default baud rate if not specified
	if baudRate == 0 {
		baudRate = 9600
	}

	m := Model{
		currentView:    ViewTerminal,
		port:           port,
		baudRate:       baudRate,
		dataBits:       8,
		stopBits:       1,
		parity:         "None",
		lineEnding:     "CRLF",
		terminalBuffer: make([]string, 0),
		maxLines:       1000,
		input:          ti,
		menuIndex:      0,
		portIndex:      0,
		connected:      false,
		ctrlAPressed:   false,
		keys:           defaultKeys,
		transferState:  NewTransferState(),
		loggingState:   NewLoggingState(),
		zone:           bubblezone.New(),
	}

	if m.zone != nil {
		m.zone.SetEnabled(true)
	}

	m.menuItems = []menuItem{
		{key: "O", label: "Serial pOrt Settings"},
		{key: "S", label: "Send files"},
		{key: "R", label: "Receive files"},
		{key: "C", label: "Clear Screen"},
		{key: "L", label: "capture to fiLe (Log)"},
		{key: "Z", label: "Help (show this menu)"},
		{key: "X", label: "eXit and reset"},
		{key: "Q", label: "Quit with no reset"},
	}

	m.helpMenuItems = []menuItem{
		// {key: "🖱️", label: "Click [Connect]/[Disconnect] button"},
		{key: "^A Z", label: "Show Help Menu"},
		{key: "^A O", label: "Serial Port Settings"},
		{key: "^A C", label: "Clear Screen"},
		{key: "^A S", label: "Send Files"},
		{key: "^A R", label: "Receive Files"},
		{key: "^A L", label: "Toggle Logging"},
		{key: "^A Q", label: "Quit"},
		{key: "^A X", label: "Exit"},
	}

	return m
}

// autoConnectMsg is sent to trigger auto-connection
type autoConnectMsg struct{}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tickCmd(),
		m.refreshPorts(),
	}

	// Auto-connect if port was specified via command line
	if m.port != "" {
		cmds = append(cmds, func() tea.Msg {
			return autoConnectMsg{}
		})
	}

	return tea.Batch(cmds...)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) refreshPorts() tea.Cmd {
	return func() tea.Msg {
		ports, err := serial.GetPortsList()
		if err != nil {
			return nil
		}
		return portsListMsg(ports)
	}
}

// readSerialCmd reads data from the serial port channel
func (m Model) readSerialCmd() tea.Cmd {
	return func() tea.Msg {
		if m.serialChan == nil {
			return nil
		}

		select {
		case data, ok := <-m.serialChan:
			if !ok {
				return nil
			}
			return serialDataMsg(data)
		default:
			return nil
		}
	}
}

// waitForSerialData waits for data from serial channel
func (m Model) waitForSerialData() tea.Cmd {
	return func() tea.Msg {
		if m.serialChan == nil {
			return nil
		}

		data, ok := <-m.serialChan
		if !ok {
			return nil
		}
		return serialDataMsg(data)
	}
}

// startSerialReader starts a goroutine to continuously read from serial port
func (m *Model) startSerialReader() {
	if m.reading {
		return
	}

	reader := m.currentSerialReader()
	if reader == nil {
		return
	}

	m.serialChan = make(chan string, 100)
	m.stopReadChan = make(chan struct{})
	m.reading = true

	go func(r io.Reader) {
		buf := make([]byte, 2048)
		for {
			select {
			case <-m.stopReadChan:
				return
			default:
			}

			n, err := r.Read(buf)
			if n > 0 {
				payload := string(buf[:n])
				select {
				case m.serialChan <- payload:
				case <-m.stopReadChan:
					return
				}
			}
			if err != nil {
				if err == io.EOF {
					time.Sleep(50 * time.Millisecond)
				}
				continue
			}
		}
	}(reader)
}

// stopSerialReader stops the serial reader goroutine
func (m *Model) stopSerialReader() {
	if m.stopReadChan != nil {
		close(m.stopReadChan)
		m.stopReadChan = nil
	}
	if m.serialChan != nil {
		close(m.serialChan)
		m.serialChan = nil
	}
	m.reading = false
}

// serialErrorMsg for serial read errors
type serialErrorMsg struct {
	err error
}

type portsListMsg []string

func (m *Model) connect() {
	if m.port == "" {
		m.lastError = "No port selected"
		return
	}

	mode := &serial.Mode{
		BaudRate: m.baudRate,
		DataBits: m.dataBits,
		StopBits: serial.StopBits(m.stopBits),
	}

	switch m.parity {
	case "Even":
		mode.Parity = serial.EvenParity
	case "Odd":
		mode.Parity = serial.OddParity
	case "Mark":
		mode.Parity = serial.MarkParity
	case "Space":
		mode.Parity = serial.SpaceParity
	default:
		mode.Parity = serial.NoParity
	}

	port, err := serial.Open(m.port, mode)
	if err != nil {
		m.lastError = err.Error()
		return
	}

	// Set read timeout to prevent blocking forever
	port.SetReadTimeout(50 * time.Millisecond)

	m.serialPort = port
	m.connected = true
	m.lastError = ""
	m.appendToBuffer("--- Connected to " + m.port + " ---")
	m.setupTrzszBridge()

	// Start the serial reader goroutine
	m.startSerialReader()
}

func (m *Model) disconnect() {
	// Stop the serial reader first
	m.stopSerialReader()
	m.teardownTrzszBridge()

	if m.serialPort != nil {
		_ = m.serialPort.Close()
		m.serialPort = nil
	}
	m.connected = false
	m.serialBuffer = "" // Clear any pending incomplete data
	m.appendToBuffer("--- Disconnected ---")
}

func (m *Model) appendToBuffer(line string) {
	m.terminalBuffer = append(m.terminalBuffer, line)
	if len(m.terminalBuffer) > m.maxLines {
		m.terminalBuffer = m.terminalBuffer[1:]
	}

	// Log the line if logging is enabled
	if m.loggingState.IsEnabled() {
		_ = m.loggingState.WriteLine(line)
	}
}

// Run starts the TUI
func Run(port string, baudRate int) error {
	m := NewModel(port, baudRate)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

func (m *Model) currentSerialReader() io.Reader {
	if m.trzFilter != nil && m.filterOutputReader != nil {
		return m.filterOutputReader
	}
	return m.serialPort
}

func (m Model) writeSerial(data []byte) {
	if len(data) == 0 {
		return
	}
	if m.filterInputWriter != nil {
		_, _ = m.filterInputWriter.Write(data)
		return
	}
	if m.serialPort != nil {
		_, _ = m.serialPort.Write(data)
	}
}

func (m *Model) setupTrzszBridge() {
	if m.trzFilter != nil || m.serialPort == nil {
		return
	}
	clientInReader, clientInWriter := io.Pipe()
	clientOutReader, clientOutWriter := io.Pipe()
	columns := m.currentTerminalColumns()
	m.trzFilter = trzsz.NewTrzszFilter(
		clientInReader,
		clientOutWriter,
		serialWriteCloser{port: m.serialPort},
		m.serialPort,
		trzsz.TrzszOptions{
			TerminalColumns: columns,
			EnableZmodem:    true,
		},
	)
	m.filterInputWriter = clientInWriter
	m.filterOutputReader = clientOutReader
}

func (m *Model) teardownTrzszBridge() {
	if m.trzFilter != nil {
		m.trzFilter.StopTransferringFiles(true)
		m.trzFilter.ResetTerminal()
		m.trzFilter = nil
	}
	if m.filterInputWriter != nil {
		_ = m.filterInputWriter.Close()
		m.filterInputWriter = nil
	}
	if m.filterOutputReader != nil {
		_ = m.filterOutputReader.Close()
		m.filterOutputReader = nil
	}
	m.clearTransferWait()
}

func (m Model) currentTerminalColumns() int32 {
	if m.width <= 0 {
		return 80
	}
	return int32(m.width)
}

type serialWriteCloser struct {
	port serial.Port
}

func (s serialWriteCloser) Write(p []byte) (int, error) {
	if s.port == nil {
		return 0, io.ErrClosedPipe
	}
	return s.port.Write(p)
}

func (serialWriteCloser) Close() error {
	return nil
}

func (m *Model) beginTransferWait(direction TransferDirection, label string) {
	m.transferPending = true
	m.transferPendingDirection = direction
	m.transferPendingLabel = label
	m.transferSpinnerIndex = 0
}

func (m *Model) clearTransferWait() {
	m.transferPending = false
	m.transferPendingLabel = ""
	m.transferSpinnerIndex = 0
}

func (m *Model) advanceTransferSpinner() {
	if !m.transferPending || len(transferSpinnerFrames) == 0 {
		return
	}
	m.transferSpinnerIndex = (m.transferSpinnerIndex + 1) % len(transferSpinnerFrames)
}

func (m Model) currentTransferStatus() string {
	if !m.transferPending || len(transferSpinnerFrames) == 0 {
		return ""
	}
	frame := transferSpinnerFrames[m.transferSpinnerIndex%len(transferSpinnerFrames)]
	label := m.transferPendingLabel
	if label == "" {
		label = transferDirectionLabel(m.transferPendingDirection)
	}
	return frame + " " + label
}
