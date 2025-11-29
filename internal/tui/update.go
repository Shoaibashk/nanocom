package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type transferResultMsg struct {
	direction TransferDirection
	path      string
	err       error
}

// Update handles Bubble Tea messages and updates the model state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reset dialog position when window resizes
		m.dialogX = 0
		m.dialogY = 0
		if m.trzFilter != nil {
			m.trzFilter.SetTerminalColumns(m.currentTerminalColumns())
		}
		return m, nil

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	case tickMsg:
		if m.transferPending {
			m.advanceTransferSpinner()
		}
		// Check for serial data if connected
		if m.connected && m.reading {
			return m, tea.Batch(tickCmd(), m.readSerialCmd())
		}
		return m, tickCmd()

	case portsListMsg:
		m.availablePorts = msg
		return m, nil

	case autoConnectMsg:
		// Auto-connect when port was specified via command line
		if m.port != "" && !m.connected {
			m.connect()
			if m.connected && m.reading {
				return m, m.waitForSerialData()
			}
		}
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

	case transferResultMsg:
		m.handleTransferResult(msg)
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
			} else if m.loggingState.IsEnabled() {
				m.appendToBuffer("--- Logging started: " + m.loggingState.GetFilePath() + " ---")
			} else {
				lines, bytes := m.loggingState.GetStats()
				m.appendToBuffer(fmt.Sprintf("--- Logging stopped (lines: %d, bytes: %d) ---", lines, bytes))
			}
			m.currentView = ViewTerminal
			return m, nil, true
		case "q", "x":
			return m, tea.Quit, true
		case "a":
			// Send Ctrl+A literally
			if m.connected {
				m.writeSerial([]byte{0x01})
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

// handleMouseEvent handles mouse events for dragging dialog windows and button clicks
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Handle button clicks in terminal view
	if m.currentView == ViewTerminal {
		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			if m.zone != nil {
				if zone := m.zone.Get(statusButtonZoneID); zone != nil && zone.InBounds(msg) {
					if m.connected {
						m.disconnect()
					} else {
						m.connect()
						if m.connected && m.reading {
							return m, m.waitForSerialData()
						}
					}
				}
			}
		}
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
		if m.connected && m.input.Value() != "" {
			data := m.input.Value() + m.getLineEndingBytes()
			m.writeSerial([]byte(data))
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
	if !m.connected || m.serialPort == nil {
		m.transferState.errorMsg = "Connect to a serial port first"
		return m, nil
	}
	if m.trzFilter == nil {
		m.transferState.errorMsg = "File transfer bridge unavailable"
		return m, nil
	}
	if m.trzFilter.IsTransferringFiles() {
		m.transferState.errorMsg = "Another transfer is already running"
		return m, nil
	}

	localPath := strings.TrimSpace(m.transferState.fileInput.Value())
	if localPath == "" {
		m.transferState.errorMsg = "Local file path is required"
		return m, nil
	}

	absPath, err := expandPath(localPath)
	if err != nil {
		m.transferState.errorMsg = fmt.Sprintf("Invalid path: %v", err)
		return m, nil
	}
	info, statErr := os.Stat(absPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			m.transferState.errorMsg = "File or directory does not exist"
		} else {
			m.transferState.errorMsg = fmt.Sprintf("Unable to access path: %v", statErr)
		}
		return m, nil
	}

	resultChan, err := m.trzFilter.OneTimeUpload([]string{absPath})
	if err != nil {
		m.transferState.errorMsg = fmt.Sprintf("Failed to start transfer: %v", err)
		return m, nil
	}

	uploadRoot := absPath
	if info != nil && !info.IsDir() {
		uploadRoot = filepath.Dir(absPath)
	}
	if uploadRoot != "" {
		m.trzFilter.SetDefaultUploadPath(uploadRoot)
	}

	protocolCmd := protocolCommand(TransferSend, m.transferState.protocol)
	protoName := m.transferState.GetProtocolName()
	message := fmt.Sprintf("--- Ready to send %s via %s ---", filepath.Base(absPath), protoName)
	if protocolCmd != "" {
		message += fmt.Sprintf(" Run '%s' on the remote device to begin.", protocolCmd)
	}
	m.appendToBuffer(message)
	m.beginTransferWait(TransferSend, fmt.Sprintf("Waiting on remote: %s", filepath.Base(absPath)))

	transferPath := absPath
	m.transferState.Reset()
	m.currentView = ViewTerminal

	return m, waitForTransferResult(resultChan, TransferSend, transferPath)
}

func (m *Model) executeReceiveFile() (tea.Model, tea.Cmd) {
	if !m.connected || m.serialPort == nil {
		m.transferState.errorMsg = "Connect to a serial port first"
		return m, nil
	}
	if m.trzFilter == nil {
		m.transferState.errorMsg = "File transfer bridge unavailable"
		return m, nil
	}
	if m.trzFilter.IsTransferringFiles() {
		m.transferState.errorMsg = "Another transfer is already running"
		return m, nil
	}

	localPath := strings.TrimSpace(m.transferState.fileInput.Value())
	if localPath == "" {
		m.transferState.errorMsg = "Destination directory is required"
		return m, nil
	}
	absPath, err := expandPath(localPath)
	if err != nil {
		m.transferState.errorMsg = fmt.Sprintf("Invalid path: %v", err)
		return m, nil
	}
	info, statErr := os.Stat(absPath)
	saveDir := absPath
	if statErr != nil {
		if os.IsNotExist(statErr) {
			if mkErr := os.MkdirAll(absPath, 0o755); mkErr != nil {
				m.transferState.errorMsg = fmt.Sprintf("Failed to create directory: %v", mkErr)
				return m, nil
			}
		} else {
			m.transferState.errorMsg = fmt.Sprintf("Unable to access path: %v", statErr)
			return m, nil
		}
	} else if info != nil && !info.IsDir() {
		saveDir = filepath.Dir(absPath)
	}
	if saveDir == "" {
		m.transferState.errorMsg = "Destination directory is required"
		return m, nil
	}

	m.trzFilter.SetDefaultDownloadPath(saveDir)
	protoName := m.transferState.GetProtocolName()
	protocolCmd := protocolCommand(TransferReceive, m.transferState.protocol)
	message := fmt.Sprintf("--- Ready to receive files into %s via %s ---", saveDir, protoName)
	if protocolCmd != "" {
		message += fmt.Sprintf(" Ask the remote device to run '%s' with the files you need.", protocolCmd)
	}
	m.appendToBuffer(message)

	m.transferState.Reset()
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
		} else if m.loggingState.IsEnabled() {
			m.appendToBuffer("--- Logging started: " + m.loggingState.GetFilePath() + " ---")
		} else {
			lines, bytes := m.loggingState.GetStats()
			m.appendToBuffer(fmt.Sprintf("--- Logging stopped (lines: %d, bytes: %d) ---", lines, bytes))
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

func (m *Model) handleTransferResult(msg transferResultMsg) {
	m.clearTransferWait()
	label := transferDirectionLabel(msg.direction)
	if msg.err != nil {
		m.appendToBuffer(fmt.Sprintf("--- %s transfer failed: %v ---", label, msg.err))
		return
	}
	if msg.direction == TransferSend {
		m.appendToBuffer(fmt.Sprintf("--- Send complete: %s ---", filepath.Base(msg.path)))
		return
	}
	m.appendToBuffer(fmt.Sprintf("--- Receive complete: %s ---", msg.path))
}

func waitForTransferResult(ch <-chan error, direction TransferDirection, path string) tea.Cmd {
	return func() tea.Msg {
		if ch == nil {
			return transferResultMsg{direction: direction, path: path, err: fmt.Errorf("transfer channel unavailable")}
		}
		err, ok := <-ch
		if !ok {
			return transferResultMsg{direction: direction, path: path}
		}
		return transferResultMsg{direction: direction, path: path, err: err}
	}
}

func transferDirectionLabel(direction TransferDirection) string {
	if direction == TransferSend {
		return "send"
	}
	return "receive"
}
