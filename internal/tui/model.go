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
	port     string
	baudRate int
	dataBits int
	stopBits int
	parity   string

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

	// Key bindings
	keys keyMap
}

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
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

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
		terminalBuffer: make([]string, 0),
		maxLines:       1000,
		input:          ti,
		menuIndex:      0,
		portIndex:      0,
		connected:      false,
		ctrlAPressed:   false,
		keys:           defaultKeys,
	}

	m.menuItems = []menuItem{
		{key: "O", label: "cOnfigure Minicom"},
		{key: "S", label: "Send files"},
		{key: "R", label: "Receive files"},
		{key: "C", label: "Clear Screen"},
		{key: "P", label: "communication Parameters"},
		{key: "L", label: "capture to fiLe (Log)"},
		{key: "Z", label: "Help (show this menu)"},
		{key: "X", label: "eXit and reset"},
		{key: "Q", label: "Quit with no reset"},
	}

	m.helpMenuItems = []menuItem{
		{key: "^A Z", label: "Show Help Menu"},
		{key: "^A P", label: "Communication Parameters"},
		{key: "^A C", label: "Clear Screen"},
		{key: "^A O", label: "Configuration Menu"},
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

type portsListMsg []string

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		return m, tickCmd()

	case portsListMsg:
		m.availablePorts = msg
		return m, nil

	case serialDataMsg:
		m.appendToBuffer(string(msg))
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	// Update text input if in terminal view
	if m.currentView == ViewTerminal && !m.ctrlAPressed {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	// Handle Ctrl+A command mode
	if m.ctrlAPressed {
		m.ctrlAPressed = false
		switch strings.ToLower(keyStr) {
		case "z":
			m.currentView = ViewHelp
			m.menuIndex = 0
			return m, nil
		case "p":
			m.currentView = ViewSerialSetup
			m.menuIndex = 0
			return m, nil
		case "c":
			m.terminalBuffer = make([]string, 0)
			m.currentView = ViewTerminal
			return m, nil
		case "o":
			m.currentView = ViewMainMenu
			m.menuIndex = 0
			return m, nil
		case "q", "x":
			return m, tea.Quit
		case "a":
			// Send Ctrl+A literally
			if m.connected && m.serialPort != nil {
				_, _ = m.serialPort.Write([]byte{0x01})
			}
			m.currentView = ViewTerminal
			return m, nil
		}
		m.currentView = ViewTerminal
		return m, nil
	}

	// Check for Ctrl+A
	if keyStr == "ctrl+a" {
		m.ctrlAPressed = true
		return m, nil
	}

	// Handle based on current view
	switch m.currentView {
	case ViewTerminal:
		return m.handleTerminalKey(msg)
	case ViewMainMenu:
		return m.handleMenuKey(msg)
	case ViewSerialSetup:
		return m.handleSerialSetupKey(msg)
	case ViewHelp:
		return m.handleHelpKey(msg)
	case ViewPortList:
		return m.handlePortListKey(msg)
	}

	return m, nil
}

func (m Model) handleTerminalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		// Send input to serial port
		if m.connected && m.serialPort != nil && m.input.Value() != "" {
			data := m.input.Value() + "\r\n"
			_, _ = m.serialPort.Write([]byte(data))
			m.appendToBuffer("> " + m.input.Value())
			m.input.SetValue("")
		}
		return m, nil
	}
	return m, nil
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
		if m.menuIndex < 4 {
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

func (m Model) selectMenuItem() (tea.Model, tea.Cmd) {
	if m.menuIndex >= len(m.menuItems) {
		return m, nil
	}

	item := m.menuItems[m.menuIndex]
	switch item.key {
	case "O":
		m.currentView = ViewSerialSetup
		m.menuIndex = 0
	case "C":
		m.terminalBuffer = make([]string, 0)
		m.currentView = ViewTerminal
	case "P":
		m.currentView = ViewSerialSetup
		m.menuIndex = 0
	case "Z":
		m.currentView = ViewHelp
		m.menuIndex = 0
	case "X", "Q":
		return m, tea.Quit
	default:
		m.currentView = ViewTerminal
	}
	return m, nil
}

func (m *Model) handleSerialSetupSelect() (tea.Model, tea.Cmd) {
	switch m.menuIndex {
	case 0: // Port selection
		m.currentView = ViewPortList
		m.portIndex = 0
		// Find current port in list
		for i, p := range m.availablePorts {
			if p == m.port {
				m.portIndex = i
				break
			}
		}
	case 4: // Connect/Disconnect
		if m.connected {
			m.disconnect()
		} else {
			m.connect()
		}
	}
	return m, nil
}

func (m *Model) adjustSerialSetting(delta int) {
	baudRates := []int{300, 1200, 2400, 4800, 9600, 19200, 38400, 57600, 115200, 230400}
	dataBitsOptions := []int{5, 6, 7, 8}
	stopBitsOptions := []int{1, 2}
	parityOptions := []string{"None", "Even", "Odd", "Mark", "Space"}

	switch m.menuIndex {
	case 1: // Baud rate
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
	case 2: // Data bits
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
	case 3: // Stop bits
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
	case 4: // Parity (reusing index 4, but we can skip for action)
		idx := 0
		for i, p := range parityOptions {
			if p == m.parity {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(parityOptions) - 1
		} else if idx >= len(parityOptions) {
			idx = 0
		}
		m.parity = parityOptions[idx]
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

	m.serialPort = port
	m.connected = true
	m.lastError = ""
	m.appendToBuffer("--- Connected to " + m.port + " ---")
}

func (m *Model) disconnect() {
	if m.serialPort != nil {
		_ = m.serialPort.Close()
		m.serialPort = nil
	}
	m.connected = false
	m.appendToBuffer("--- Disconnected ---")
}

func (m *Model) appendToBuffer(line string) {
	m.terminalBuffer = append(m.terminalBuffer, line)
	if len(m.terminalBuffer) > m.maxLines {
		m.terminalBuffer = m.terminalBuffer[1:]
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
	}

	return content
}

func (m Model) renderTerminalView() string {
	// Title bar
	title := " nanocom 1.0 - Terminal "
	titleBar := titleBarStyle.Width(m.width).Render(title)

	// Calculate terminal height
	termHeight := m.height - 4 // title + status + input

	// Terminal content
	var termContent strings.Builder
	startIdx := 0
	if len(m.terminalBuffer) > termHeight {
		startIdx = len(m.terminalBuffer) - termHeight
	}

	for i := startIdx; i < len(m.terminalBuffer); i++ {
		termContent.WriteString(m.terminalBuffer[i])
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

	// Input line
	inputLine := inputStyle.Width(m.width).Render(m.input.View())

	// Status bar
	status := m.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left, titleBar, terminal, inputLine, status)
}

func (m Model) renderStatusBar() string {
	portInfo := fmt.Sprintf(" %s | %d %d%s%d ",
		m.port, m.baudRate, m.dataBits, string(m.parity[0]), m.stopBits)

	connStatus := "OFFLINE"
	if m.connected {
		connStatus = "ONLINE"
	}

	timeStr := time.Now().Format("15:04:05")
	hint := "CTRL-A Z for help"

	left := portInfo + " | " + connStatus
	right := hint + " | " + timeStr

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	statusContent := left + strings.Repeat(" ", gap) + right
	return statusBarStyle.Width(m.width).Render(statusContent)
}

func (m Model) renderMainMenu() string {
	// Background terminal view (dimmed)
	bg := m.renderTerminalView()

	// Menu overlay
	var menuContent strings.Builder
	menuContent.WriteString(menuTitleStyle.Render("nanocom Command Summary"))
	menuContent.WriteString("\n\n")

	for i, item := range m.menuItems {
		hotkey := menuHotkeyStyle.Render(item.key)
		label := " " + item.label

		if i == m.menuIndex {
			line := menuSelectedStyle.Render(hotkey + label)
			menuContent.WriteString(line)
		} else {
			line := menuItemStyle.Render(hotkey + label)
			menuContent.WriteString(line)
		}
		menuContent.WriteString("\n")
	}

	menuContent.WriteString("\n")
	menuContent.WriteString(helpStyle.Render("Use ↑↓ to navigate, Enter to select, Esc to close"))

	menu := menuStyle.Render(menuContent.String())

	// Center the menu
	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)
	x := (m.width - menuWidth) / 2
	y := (m.height - menuHeight) / 2

	return m.overlayAt(bg, menu, x, y)
}

func (m Model) renderSerialSetup() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("Serial Port Setup"))
	content.WriteString("\n\n")

	settings := []struct {
		label string
		value string
	}{
		{"Serial Device", m.port},
		{"Bps/Par/Bits", fmt.Sprintf("%d %s %d", m.baudRate, m.parity, m.dataBits)},
		{"Data Bits", fmt.Sprintf("%d", m.dataBits)},
		{"Stop Bits", fmt.Sprintf("%d", m.stopBits)},
		{m.getConnectLabel(), ""},
	}

	for i, s := range settings {
		prefix := "  "
		if i == m.menuIndex {
			prefix = "> "
		}

		label := labelStyle.Render(s.label + ":")
		value := s.value

		if i == m.menuIndex {
			line := menuSelectedStyle.Render(fmt.Sprintf("%s%-20s %s", prefix, s.label+":", value))
			content.WriteString(line)
		} else {
			content.WriteString(fmt.Sprintf("%s%s %s", prefix, label, value))
		}
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("↑↓: navigate, ←→: change value, Enter: select, Esc: back"))

	if m.lastError != "" {
		content.WriteString("\n")
		content.WriteString(errorStyle.Render("Error: " + m.lastError))
	}

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)
	x := (m.width - menuWidth) / 2
	y := (m.height - menuHeight) / 2

	return m.overlayAt(bg, menu, x, y)
}

func (m Model) getConnectLabel() string {
	if m.connected {
		return "Disconnect"
	}
	return "Connect"
}

func (m Model) renderHelpMenu() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("nanocom Help"))
	content.WriteString("\n\n")

	for _, item := range m.helpMenuItems {
		hotkey := menuHotkeyStyle.Render(fmt.Sprintf("%-8s", item.key))
		label := item.label
		content.WriteString(fmt.Sprintf(" %s  %s\n", hotkey, label))
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("Press any key to close"))

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)
	x := (m.width - menuWidth) / 2
	y := (m.height - menuHeight) / 2

	return m.overlayAt(bg, menu, x, y)
}

func (m Model) renderPortList() string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render("Select Serial Port"))
	content.WriteString("\n\n")

	if len(m.availablePorts) == 0 {
		content.WriteString(errorStyle.Render("No serial ports found"))
	} else {
		for i, port := range m.availablePorts {
			if i == m.portIndex {
				content.WriteString(menuSelectedStyle.Render("> " + port))
			} else {
				content.WriteString(menuItemStyle.Render("  " + port))
			}
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("↑↓: navigate, Enter: select, Esc: back"))

	menu := menuStyle.Render(content.String())

	menuWidth := lipgloss.Width(menu)
	menuHeight := lipgloss.Height(menu)
	x := (m.width - menuWidth) / 2
	y := (m.height - menuHeight) / 2

	return m.overlayAt(bg, menu, x, y)
}

func (m Model) overlayAt(bg, overlay string, x, y int) string {
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
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
