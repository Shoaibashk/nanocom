/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// tickMsg for periodic updates (like time)
type tickMsg time.Time

// serialDataMsg for incoming serial data
type serialDataMsg string

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

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		m.refreshPorts(),
	)
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
	if m.serialPort == nil || m.reading {
		return
	}

	m.serialChan = make(chan string, 100)
	m.stopReadChan = make(chan struct{})
	m.reading = true

	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-m.stopReadChan:
				return
			default:
				n, err := m.serialPort.Read(buf)
				if err != nil {
					continue // Timeout or error, keep trying
				}
				if n > 0 {
					select {
					case m.serialChan <- string(buf[:n]):
					case <-m.stopReadChan:
						return
					}
				}
			}
		}
	}()
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

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reset dialog position when window resizes
		m.dialogX = 0
		m.dialogY = 0
		return m, nil

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	case tickMsg:
		// Check for serial data if connected
		if m.connected && m.reading {
			return m, tea.Batch(tickCmd(), m.readSerialCmd())
		}
		return m, tickCmd()

	case portsListMsg:
		m.availablePorts = msg
		return m, nil

	case serialDataMsg:
		// Append received data to serial buffer
		m.serialBuffer += string(msg)

		// Process complete lines (ending with \n or \r)
		for {
			// Find the first line terminator (\n or \r)
			idxN := strings.Index(m.serialBuffer, "\n")
			idxR := strings.Index(m.serialBuffer, "\r")

			idx := -1
			if idxN == -1 && idxR == -1 {
				break // No complete line yet
			} else if idxN == -1 {
				idx = idxR
			} else if idxR == -1 {
				idx = idxN
			} else {
				// Both found, use the earlier one
				if idxR < idxN {
					idx = idxR
				} else {
					idx = idxN
				}
			}

			// Extract the complete line
			line := m.serialBuffer[:idx]
			m.serialBuffer = m.serialBuffer[idx+1:]

			// Skip empty lines and handle \r\n sequence
			if len(m.serialBuffer) > 0 && m.serialBuffer[0] == '\n' {
				m.serialBuffer = m.serialBuffer[1:] // Skip the \n after \r
			}

			// Clean up any remaining carriage returns
			line = strings.Trim(line, "\r\n")

			// Only add non-empty lines
			if line != "" {
				m.appendToBuffer(line)
			}
		}

		// Continue reading - wait for more data
		if m.connected && m.reading {
			return m, m.waitForSerialData()
		}
		return m, nil

	case serialErrorMsg:
		// Handle serial read error
		if m.connected {
			m.appendToBuffer("--- Serial error: " + msg.err.Error() + " ---")
		}
		return m, nil

	case tea.KeyMsg:
		// Handle special keys first
		newModel, cmd, handled := m.handleKeyPress(msg)
		if handled {
			return newModel, cmd
		}
		// If not handled and in terminal view, let the input handle it
		if m.currentView == ViewTerminal && !m.ctrlAPressed {
			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			return m, inputCmd
		}
		return newModel, cmd
	}

	// Update text input for other messages if in terminal view
	if m.currentView == ViewTerminal && !m.ctrlAPressed {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	keyStr := msg.String()

	// Handle Ctrl+A command mode
	if m.ctrlAPressed {
		m.ctrlAPressed = false
		switch strings.ToLower(keyStr) {
		case "z":
			m.currentView = ViewHelp
			m.menuIndex = 0
			m.dialogX = 0
			m.dialogY = 0
			return m, nil, true
		case "c":
			m.terminalBuffer = make([]string, 0)
			m.currentView = ViewTerminal
			return m, nil, true
		case "o":
			m.currentView = ViewMainMenu
			m.menuIndex = 0
			m.dialogX = 0
			m.dialogY = 0
			return m, nil, true
		case "s":
			// Send Files
			m.transferState.Reset()
			m.transferState.direction = TransferSend
			m.transferState.fileInput.Focus()
			m.currentView = ViewSendFile
			m.dialogX = 0
			m.dialogY = 0
			return m, nil, true
		case "r":
			// Receive Files
			m.transferState.Reset()
			m.transferState.direction = TransferReceive
			m.transferState.fileInput.Focus()
			m.currentView = ViewReceiveFile
			m.dialogX = 0
			m.dialogY = 0
			return m, nil, true
		case "l":
			// Toggle Logging
			err := m.loggingState.Toggle("")
			if err != nil {
				m.lastError = err.Error()
			} else {
				if m.loggingState.IsEnabled() {
					m.appendToBuffer("--- Logging started: " + m.loggingState.GetFilePath() + " ---")
				} else {
					lines, bytes := m.loggingState.GetStats()
					m.appendToBuffer(fmt.Sprintf("--- Logging stopped (lines: %d, bytes: %d) ---", lines, bytes))
				}
			}
			m.currentView = ViewTerminal
			return m, nil, true
		case "q", "x":
			return m, tea.Quit, true
		case "a":
			// Send Ctrl+A literally
			if m.connected && m.serialPort != nil {
				_, _ = m.serialPort.Write([]byte{0x01})
			}
			m.currentView = ViewTerminal
			return m, nil, true
		}
		m.currentView = ViewTerminal
		return m, nil, true
	}

	// Check for Ctrl+A
	if keyStr == "ctrl+a" {
		m.ctrlAPressed = true
		return m, nil, true
	}

	// Handle based on current view
	switch m.currentView {
	case ViewTerminal:
		return m.handleTerminalKey(msg)
	case ViewMainMenu:
		model, cmd := m.handleMenuKey(msg)
		return model, cmd, true
	case ViewSerialSetup:
		model, cmd := m.handleSerialSetupKey(msg)
		return model, cmd, true
	case ViewHelp:
		model, cmd := m.handleHelpKey(msg)
		return model, cmd, true
	case ViewPortList:
		model, cmd := m.handlePortListKey(msg)
		return model, cmd, true
	case ViewSendFile:
		model, cmd := m.handleSendFileKey(msg)
		return model, cmd, true
	case ViewReceiveFile:
		model, cmd := m.handleReceiveFileKey(msg)
		return model, cmd, true
	}

	return m, nil, false
}

// handleMouseEvent handles mouse events for dragging dialog windows
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only handle mouse events when a dialog is visible
	if m.currentView == ViewTerminal {
		return m, nil
	}

	// Calculate estimated dialog dimensions based on view type
	dialogWidth, dialogHeight := m.getEstimatedDialogSize()
	m.dialogWidth = dialogWidth
	m.dialogHeight = dialogHeight

	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button == tea.MouseButtonLeft {
			// Check if click is within dialog title bar area (first 3 rows of dialog)
			dialogX := m.getDialogX()
			dialogY := m.getDialogY()
			if msg.X >= dialogX && msg.X < dialogX+dialogWidth &&
				msg.Y >= dialogY && msg.Y < dialogY+3 {
				m.isDragging = true
				m.dragOffsetX = msg.X - dialogX
				m.dragOffsetY = msg.Y - dialogY
			}
		}
	case tea.MouseActionRelease:
		m.isDragging = false
	case tea.MouseActionMotion:
		if m.isDragging {
			newX := msg.X - m.dragOffsetX
			newY := msg.Y - m.dragOffsetY

			// Clamp to viewport bounds
			if newX < 0 {
				newX = 0
			}
			if newY < 0 {
				newY = 0
			}
			if newX+dialogWidth > m.width {
				newX = m.width - dialogWidth
			}
			if newY+dialogHeight > m.height {
				newY = m.height - dialogHeight
			}

			m.dialogX = newX
			m.dialogY = newY
		}
	}

	return m, nil
}

// getEstimatedDialogSize returns estimated dialog dimensions based on current view
func (m Model) getEstimatedDialogSize() (int, int) {
	switch m.currentView {
	case ViewMainMenu:
		return 40, len(m.menuItems) + 8
	case ViewSerialSetup:
		return 50, 12
	case ViewHelp:
		return 45, len(m.helpMenuItems) + 8
	case ViewPortList:
		ports := len(m.availablePorts)
		if ports == 0 {
			ports = 1
		}
		return 40, ports + 8
	case ViewSendFile, ViewReceiveFile:
		return 60, 14
	default:
		return 40, 15
	}
}

// getDialogX returns the current dialog X position
func (m Model) getDialogX() int {
	if m.dialogX == 0 {
		width, _ := m.getEstimatedDialogSize()
		return (m.width - width) / 2
	}
	return m.dialogX
}

// getDialogY returns the current dialog Y position
func (m Model) getDialogY() int {
	if m.dialogY == 0 {
		_, height := m.getEstimatedDialogSize()
		return (m.height - height) / 2
	}
	return m.dialogY
}

func (m Model) handleTerminalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit, true
	case "enter":
		// Send input to serial port
		if m.connected && m.serialPort != nil && m.input.Value() != "" {
			data := m.input.Value() + m.getLineEndingBytes()
			_, _ = m.serialPort.Write([]byte(data))
			m.appendToBuffer("> " + m.input.Value())
			m.input.SetValue("")
		} else if m.input.Value() != "" {
			// Echo locally even when not connected
			m.appendToBuffer("> " + m.input.Value())
			m.input.SetValue("")
		}
		return m, nil, true
	}
	// Return false to let the textinput handle other keys
	return m, nil, false
}

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
	case "enter":
		return m.selectMenuItem()
	case "esc":
		m.currentView = ViewTerminal
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleSerialSetupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < serialSetupMenuItemCount-1 {
			m.menuIndex++
		}
	case "enter":
		return m.handleSerialSetupSelect()
	case "esc":
		m.currentView = ViewTerminal
	case "left", "h":
		m.adjustSerialSetting(-1)
	case "right", "l":
		m.adjustSerialSetting(1)
	}
	return m, nil
}

func (m Model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", "q", "z":
		m.currentView = ViewTerminal
	}
	return m, nil
}

func (m Model) handlePortListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.portIndex > 0 {
			m.portIndex--
		}
	case "down", "j":
		if m.portIndex < len(m.availablePorts)-1 {
			m.portIndex++
		}
	case "enter":
		if len(m.availablePorts) > 0 {
			m.port = m.availablePorts[m.portIndex]
		}
		m.currentView = ViewSerialSetup
	case "esc":
		m.currentView = ViewSerialSetup
	}
	return m, nil
}

func (m Model) handleSendFileKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.transferState.Reset()
		m.currentView = ViewTerminal
		return m, nil
	case "tab":
		m.transferState.FocusNext()
		return m, nil
	case "shift+tab":
		m.transferState.FocusPrev()
		return m, nil
	case "left", "right":
		if m.transferState.inputFocus == -1 {
			// On protocol selector
			delta := 1
			if msg.String() == "left" {
				delta = -1
			}
			m.transferState.CycleProtocol(delta)
		}
		return m, nil
	case "enter":
		return m.executeSendFile()
	case "ctrl+p":
		// Toggle protocol with Ctrl+P
		m.transferState.CycleProtocol(1)
		return m, nil
	}

	// Update the focused input field
	var cmd tea.Cmd
	m.transferState.fileInput, cmd = m.transferState.fileInput.Update(msg)
	return m, cmd
}

func (m Model) handleReceiveFileKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.transferState.Reset()
		m.currentView = ViewTerminal
		return m, nil
	case "tab":
		m.transferState.FocusNext()
		return m, nil
	case "shift+tab":
		m.transferState.FocusPrev()
		return m, nil
	case "left", "right":
		if m.transferState.inputFocus == -1 {
			delta := 1
			if msg.String() == "left" {
				delta = -1
			}
			m.transferState.CycleProtocol(delta)
		}
		return m, nil
	case "enter":
		return m.executeReceiveFile()
	case "ctrl+p":
		m.transferState.CycleProtocol(1)
		return m, nil
	}

	// Update the focused input field
	var cmd tea.Cmd
	m.transferState.fileInput, cmd = m.transferState.fileInput.Update(msg)
	return m, cmd
}

func (m *Model) executeSendFile() (tea.Model, tea.Cmd) {
	localPath := m.transferState.fileInput.Value()

	if localPath == "" {
		m.transferState.errorMsg = "Local file path is required"
		return m, nil
	}

	// Protocol not yet implemented (Zmodem, Xmodem, Ymodem)
	m.transferState.errorMsg = "Protocol not yet implemented"

	m.currentView = ViewTerminal
	return m, nil
}

func (m *Model) executeReceiveFile() (tea.Model, tea.Cmd) {
	localPath := m.transferState.fileInput.Value()

	if localPath == "" {
		m.transferState.errorMsg = "Local file path is required"
		return m, nil
	}

	// Protocol not yet implemented (Zmodem, Xmodem, Ymodem)
	m.transferState.errorMsg = "Protocol not yet implemented"

	m.currentView = ViewTerminal
	return m, nil
}

func (m Model) selectMenuItem() (tea.Model, tea.Cmd) {
	if m.menuIndex >= len(m.menuItems) {
		return m, nil
	}

	item := m.menuItems[m.menuIndex]
	switch item.key {
	case "O":
		m.currentView = ViewSerialSetup
		m.menuIndex = 0
		m.dialogX = 0
		m.dialogY = 0
	case "S":
		// Send Files
		m.transferState.Reset()
		m.transferState.direction = TransferSend
		m.transferState.fileInput.Focus()
		m.currentView = ViewSendFile
		m.dialogX = 0
		m.dialogY = 0
	case "R":
		// Receive Files
		m.transferState.Reset()
		m.transferState.direction = TransferReceive
		m.transferState.fileInput.Focus()
		m.currentView = ViewReceiveFile
		m.dialogX = 0
		m.dialogY = 0
	case "L":
		// Toggle Logging
		err := m.loggingState.Toggle("")
		if err != nil {
			m.lastError = err.Error()
		} else {
			if m.loggingState.IsEnabled() {
				m.appendToBuffer("--- Logging started: " + m.loggingState.GetFilePath() + " ---")
			} else {
				lines, bytes := m.loggingState.GetStats()
				m.appendToBuffer(fmt.Sprintf("--- Logging stopped (lines: %d, bytes: %d) ---", lines, bytes))
			}
		}
		m.currentView = ViewTerminal
	case "C":
		m.terminalBuffer = make([]string, 0)
		m.currentView = ViewTerminal
	case "Z":
		m.currentView = ViewHelp
		m.menuIndex = 0
		m.dialogX = 0
		m.dialogY = 0
	case "X", "Q":
		return m, tea.Quit
	default:
		m.currentView = ViewTerminal
	}
	return m, nil
}

func (m *Model) handleSerialSetupSelect() (tea.Model, tea.Cmd) {
	switch m.menuIndex {
	case serialSetupPort:
		m.currentView = ViewPortList
		m.portIndex = 0
		m.dialogX = 0
		m.dialogY = 0
		// Find current port in list
		for i, p := range m.availablePorts {
			if p == m.port {
				m.portIndex = i
				break
			}
		}
	case serialSetupConnect:
		if m.connected {
			m.disconnect()
		} else {
			m.connect()
			// Start waiting for serial data after connecting
			if m.connected && m.reading {
				return m, m.waitForSerialData()
			}
		}
	}
	return m, nil
}

func (m *Model) adjustSerialSetting(delta int) {
	baudRates := []int{300, 1200, 2400, 4800, 9600, 19200, 38400, 57600, 115200, 230400}
	dataBitsOptions := []int{5, 6, 7, 8}
	stopBitsOptions := []int{1, 2}
	lineEndingOptions := []string{"None", "NL", "CR", "CRLF"}

	switch m.menuIndex {
	case serialSetupBaudRate:
		idx := 0
		for i, b := range baudRates {
			if b == m.baudRate {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(baudRates) - 1
		} else if idx >= len(baudRates) {
			idx = 0
		}
		m.baudRate = baudRates[idx]
	case serialSetupDataBits:
		idx := 0
		for i, d := range dataBitsOptions {
			if d == m.dataBits {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(dataBitsOptions) - 1
		} else if idx >= len(dataBitsOptions) {
			idx = 0
		}
		m.dataBits = dataBitsOptions[idx]
	case serialSetupStopBits:
		idx := 0
		for i, s := range stopBitsOptions {
			if s == m.stopBits {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(stopBitsOptions) - 1
		} else if idx >= len(stopBitsOptions) {
			idx = 0
		}
		m.stopBits = stopBitsOptions[idx]
	case serialSetupLineEnding:
		idx := 0
		for i, le := range lineEndingOptions {
			if le == m.lineEnding {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(lineEndingOptions) - 1
		} else if idx >= len(lineEndingOptions) {
			idx = 0
		}
		m.lineEnding = lineEndingOptions[idx]
	}
}

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

	// Start the serial reader goroutine
	m.startSerialReader()
}

func (m *Model) disconnect() {
	// Stop the serial reader first
	m.stopSerialReader()

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

// View renders the TUI
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content string

	switch m.currentView {
	case ViewTerminal:
		content = m.renderTerminalView()
	case ViewMainMenu:
		content = m.renderMainMenu()
	case ViewSerialSetup:
		content = m.renderSerialSetup()
	case ViewHelp:
		content = m.renderHelpMenu()
	case ViewPortList:
		content = m.renderPortList()
	case ViewSendFile:
		content = m.renderSendFileView()
	case ViewReceiveFile:
		content = m.renderReceiveFileView()
	}

	return content
}

func (m Model) renderTerminalView() string {
	// Title bar with Charm-style branding
	logoText := "◆ nanocom"
	versionText := "v1.0"
	title := logoStyle.Render(logoText) + " " + helpStyle.Render(versionText)
	titleBar := titleBarStyle.Width(m.width).Render(title)

	// Calculate terminal height
	termHeight := m.height - 5 // title + status + input + borders

	// Terminal content with styled output
	var termContent strings.Builder
	startIdx := 0
	if len(m.terminalBuffer) > termHeight {
		startIdx = len(m.terminalBuffer) - termHeight
	}

	for i := startIdx; i < len(m.terminalBuffer); i++ {
		line := m.terminalBuffer[i]
		// Style sent messages differently
		if strings.HasPrefix(line, "> ") {
			termContent.WriteString(inputPromptStyle.Render(">") + " " + line[2:])
		} else if strings.HasPrefix(line, "---") {
			termContent.WriteString(helpStyle.Render(line))
		} else {
			termContent.WriteString(line)
		}
		if i < len(m.terminalBuffer)-1 {
			termContent.WriteString("\n")
		}
	}

	// Pad empty lines
	lineCount := len(m.terminalBuffer) - startIdx
	for i := lineCount; i < termHeight; i++ {
		termContent.WriteString("\n")
	}

	terminal := terminalStyle.Width(m.width).Height(termHeight).Render(termContent.String())

	// Input line with prompt
	prompt := inputPromptStyle.Render("❯ ")
	inputField := m.input.View()
	inputLine := inputStyle.Width(m.width - 4).Render(prompt + inputField)

	// Status bar
	status := m.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left, titleBar, terminal, inputLine, status)
}

func (m Model) renderStatusBar() string {
	parityChar := "N"
	if len(m.parity) > 0 {
		parityChar = string(m.parity[0])
	}

	// Port info with styled badge
	portDisplay := m.port
	if portDisplay == "" {
		portDisplay = "No port"
	}
	portInfo := fmt.Sprintf("⚡ %s • %d %d%s%d",
		portDisplay, m.baudRate, m.dataBits, parityChar, m.stopBits)

	// Connection status badge
	var connBadge string
	if m.connected {
		connBadge = connectedStyle.Render("● ONLINE")
	} else {
		connBadge = disconnectedStyle.Render("○ OFFLINE")
	}

	// Logging status badge
	var logBadge string
	if m.loggingState.IsEnabled() {
		logBadge = warningStyle.Render(" 📝 LOG")
	}

	timeStr := time.Now().Format("15:04:05")
	hint := helpStyle.Render("Ctrl+A Z for help")

	// Build status bar with proper spacing
	left := portInfo + "  " + connBadge + logBadge
	right := hint + "  " + helpStyle.Render(timeStr)

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	gap := m.width - leftWidth - rightWidth - 4
	if gap < 0 {
		gap = 1
	}

	statusContent := left + strings.Repeat(" ", gap) + right
	return statusBarStyle.Width(m.width).Render(statusContent)
}

func (m Model) renderMainMenu() string {
	// Background terminal view (dimmed)
	bg := m.renderTerminalView()

	// Menu overlay with Charm styling
	var menuContent strings.Builder
	menuContent.WriteString(menuTitleStyle.Render("◆ nanocom Command Summary"))
	menuContent.WriteString("\n\n")

	for i, item := range m.menuItems {
		if i == m.menuIndex {
			line := menuSelectedStyle.Render(fmt.Sprintf("[%s] %s", item.key, item.label))
			menuContent.WriteString(line)
		} else {
			hotkey := menuHotkeyStyle.Render(fmt.Sprintf("[%s]", item.key))
			line := menuItemStyle.Render(fmt.Sprintf("%s %s", hotkey, item.label))
			menuContent.WriteString(line)
		}
		menuContent.WriteString("\n")
	}

	menuContent.WriteString("\n")
	menuContent.WriteString(helpStyle.Render("↑/↓ navigate • enter select • esc close"))

	menu := menuStyle.Render(menuContent.String())

	// Calculate menu dimensions and position
	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)

	// Use custom position if dragged, otherwise center
	x := m.dialogX
	y := m.dialogY
	if x == 0 && y == 0 {
		x = (m.width - menuWidth) / 2
		y = (m.height - menuHeight) / 2
	}

	return m.overlayAtWithSize(bg, menu, x, y, menuWidth, menuHeight)
}

func (m Model) renderSerialSetup() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("Serial Port Settings"))
	content.WriteString("\n\n")

	// Port display
	portDisplay := m.port
	if portDisplay == "" {
		portDisplay = "(none selected)"
	}

	settings := []struct {
		label string
		value string
	}{
		{"Serial Device", portDisplay},
		{"Baud Rate", fmt.Sprintf("%d", m.baudRate)},
		{"Data Bits", fmt.Sprintf("%d", m.dataBits)},
		{"Stop Bits", fmt.Sprintf("%d", m.stopBits)},
		{"Line Ending", m.getLineEndingDisplay()},
		{m.getConnectLabel(), ""},
	}

	// Calculate max label length for alignment
	maxLabelLen := 0
	for _, s := range settings {
		if len(s.label) > maxLabelLen {
			maxLabelLen = len(s.label)
		}
	}

	for i, s := range settings {
		// Pad the label
		paddedLabel := fmt.Sprintf("%-*s", maxLabelLen, s.label)

		if i == m.menuIndex {
			if s.value != "" {
				line := menuSelectedStyle.Render(fmt.Sprintf("%s : %s", paddedLabel, s.value))
				content.WriteString(line)
			} else {
				line := menuSelectedStyle.Render(paddedLabel)
				content.WriteString(line)
			}
		} else {
			// Style the label
			styledLabel := labelStyle.Render(paddedLabel)

			var line string
			if s.value != "" {
				line = menuItemStyle.Render(fmt.Sprintf("%s : %s", styledLabel, s.value))
			} else {
				line = menuItemStyle.Render(styledLabel)
			}
			content.WriteString(line)
		}
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("↑/↓ navigate • ←/→ change • enter select • esc back"))

	if m.lastError != "" {
		content.WriteString("\n\n")
		content.WriteString(errorStyle.Render("⚠ " + m.lastError))
	}

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)

	// Use custom position if dragged, otherwise center
	x := m.dialogX
	y := m.dialogY
	if x == 0 && y == 0 {
		x = (m.width - menuWidth) / 2
		y = (m.height - menuHeight) / 2
	}

	return m.overlayAtWithSize(bg, menu, x, y, menuWidth, menuHeight)
}

func (m Model) getConnectLabel() string {
	if m.connected {
		return "Disconnect"
	}
	return "Connect"
}

func (m Model) getLineEndingDisplay() string {
	switch m.lineEnding {
	case "None":
		return "No Line Ending"
	case "NL":
		return "Newline (\\n)"
	case "CR":
		return "Carriage Return (\\r)"
	case "CRLF":
		return "Both NL & CR (\\r\\n)"
	default:
		return m.lineEnding
	}
}

func (m Model) getLineEndingBytes() string {
	switch m.lineEnding {
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

func (m Model) renderHelpMenu() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("◆ nanocom Help"))
	content.WriteString("\n\n")

	for _, item := range m.helpMenuItems {
		hotkey := menuHotkeyStyle.Render(fmt.Sprintf("%-10s", item.key))
		label := item.label
		content.WriteString(fmt.Sprintf("  %s  %s\n", hotkey, label))
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("Press any key to close"))

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)

	// Use custom position if dragged, otherwise center
	x := m.dialogX
	y := m.dialogY
	if x == 0 && y == 0 {
		x = (m.width - menuWidth) / 2
		y = (m.height - menuHeight) / 2
	}

	return m.overlayAtWithSize(bg, menu, x, y, menuWidth, menuHeight)
}

func (m Model) renderPortList() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("Select Serial Port"))
	content.WriteString("\n\n")

	if len(m.availablePorts) == 0 {
		content.WriteString(errorStyle.Render("⚠ No serial ports found"))
	} else {
		for i, port := range m.availablePorts {
			if i == m.portIndex {
				content.WriteString(menuSelectedStyle.Render(" ● " + port + " "))
			} else {
				content.WriteString(menuItemStyle.Render("   " + port))
			}
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("↑/↓ navigate • enter select • esc back"))

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)

	// Use custom position if dragged, otherwise center
	x := m.dialogX
	y := m.dialogY
	if x == 0 && y == 0 {
		x = (m.width - menuWidth) / 2
		y = (m.height - menuHeight) / 2
	}

	return m.overlayAtWithSize(bg, menu, x, y, menuWidth, menuHeight)
}

func (m Model) renderSendFileView() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("Send Files"))
	content.WriteString("\n\n")

	// Protocol selector
	protocolLabel := labelStyle.Render("Protocol:")
	protocolValue := m.transferState.GetProtocolName()
	content.WriteString(fmt.Sprintf("  %s %s (Ctrl+P to change)\n\n", protocolLabel, protocolValue))

	// Input field - only local file
	label := labelStyle.Render("Local File:")
	content.WriteString(fmt.Sprintf("  %s %s\n", menuSelectedStyle.Render(">"), label))
	content.WriteString(fmt.Sprintf("    %s\n", m.transferState.fileInput.View()))

	content.WriteString("\n")

	// Status/error messages
	if m.transferState.errorMsg != "" {
		content.WriteString(errorStyle.Render(m.transferState.errorMsg))
		content.WriteString("\n")
	}
	if m.transferState.statusMsg != "" {
		content.WriteString(successStyle.Render(m.transferState.statusMsg))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("enter: send | esc: cancel"))

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)

	// Use custom position if dragged, otherwise center
	x := m.dialogX
	y := m.dialogY
	if x == 0 && y == 0 {
		x = (m.width - menuWidth) / 2
		y = (m.height - menuHeight) / 2
	}

	return m.overlayAtWithSize(bg, menu, x, y, menuWidth, menuHeight)
}

func (m Model) renderReceiveFileView() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("Receive Files"))
	content.WriteString("\n\n")

	// Protocol selector
	protocolLabel := labelStyle.Render("Protocol:")
	protocolValue := m.transferState.GetProtocolName()
	content.WriteString(fmt.Sprintf("  %s %s (Ctrl+P to change)\n\n", protocolLabel, protocolValue))

	// Input field - only local file destination
	label := labelStyle.Render("Local File:")
	content.WriteString(fmt.Sprintf("  %s %s\n", menuSelectedStyle.Render(">"), label))
	content.WriteString(fmt.Sprintf("    %s\n", m.transferState.fileInput.View()))

	content.WriteString("\n")

	// Status/error messages
	if m.transferState.errorMsg != "" {
		content.WriteString(errorStyle.Render(m.transferState.errorMsg))
		content.WriteString("\n")
	}
	if m.transferState.statusMsg != "" {
		content.WriteString(successStyle.Render(m.transferState.statusMsg))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("enter: receive | esc: cancel"))

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)

	// Use custom position if dragged, otherwise center
	x := m.dialogX
	y := m.dialogY
	if x == 0 && y == 0 {
		x = (m.width - menuWidth) / 2
		y = (m.height - menuHeight) / 2
	}

	return m.overlayAtWithSize(bg, menu, x, y, menuWidth, menuHeight)
}

func (m Model) overlayAt(bg, overlay string, x, y int) string {
	return m.overlayAtWithSize(bg, overlay, x, y, lipgloss.Width(overlay), lipgloss.Height(overlay))
}

func (m Model) overlayAtWithSize(bg, overlay string, x, y, width, height int) string {
	bgLines := strings.Split(bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	// Ensure we have enough lines
	for len(bgLines) < m.height {
		bgLines = append(bgLines, strings.Repeat(" ", m.width))
	}

	for i, overlayLine := range overlayLines {
		lineY := y + i
		if lineY >= 0 && lineY < len(bgLines) {
			bgLine := bgLines[lineY]

			// Pad bgLine if needed
			for len(bgLine) < x+lipgloss.Width(overlayLine) {
				bgLine += " "
			}

			// Construct the new line
			prefix := ""
			if x > 0 && x <= len(bgLine) {
				prefix = bgLine[:x]
			} else if x > len(bgLine) {
				prefix = bgLine + strings.Repeat(" ", x-len(bgLine))
			}

			suffix := ""
			endX := x + lipgloss.Width(overlayLine)
			if endX < len(bgLine) {
				suffix = bgLine[endX:]
			}

			bgLines[lineY] = prefix + overlayLine + suffix
		}
	}

	return strings.Join(bgLines, "\n")
}

// Run starts the TUI
func Run(port string, baudRate int) error {
	m := NewModel(port, baudRate)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
