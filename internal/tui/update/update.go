/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package update

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/shoaibashk/nanocom/internal/tui/menu"
	"github.com/shoaibashk/nanocom/internal/tui/model"
	"github.com/shoaibashk/nanocom/internal/tui/serial"
	"github.com/shoaibashk/nanocom/internal/tui/state"
	"github.com/shoaibashk/nanocom/internal/tui/transfer"
	"github.com/shoaibashk/nanocom/internal/tui/views"
)

const (
	// DialogHeaderHeight is the height reserved for dialog title bar (for dragging)
	DialogHeaderHeight = 3
)

// TransferResultMsg contains the result of a file transfer operation.
type TransferResultMsg struct {
	Direction transfer.Direction
	Path      string
	Err       error
}

// Update handles Bubble Tea messages and updates the model state
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.DialogX = 0
		m.DialogY = 0
		if m.TransferBridge.IsActive() {
			m.TransferBridge.SetTerminalColumns(m.CurrentTerminalColumns())
		}
		return m, nil

	case tea.MouseMsg:
		return handleMouseEvent(m, msg)

	case state.TickMsg:
		if m.TransferPending {
			m.AdvanceTransferSpinner()
		}
		if m.Serial.Connected && m.Reader.Reading {
			return m, tea.Batch(model.TickCmd(), m.ReadSerialCmd())
		}
		return m, model.TickCmd()

	case state.PortsListMsg:
		m.AvailablePorts = msg
		return m, nil

	case state.AutoConnectMsg:
		if m.Serial.Config.Port != "" && !m.Serial.Connected {
			m.Connect()
			if m.Serial.Connected && m.Reader.Reading {
				return m, m.WaitForSerialData()
			}
		}
		return m, nil

	case state.SerialDataMsg:
		processSerialData(&m, string(msg))

		if m.Serial.Connected && m.Reader.Reading {
			return m, m.WaitForSerialData()
		}
		return m, nil

	case state.SerialErrorMsg:
		if m.Serial.Connected {
			m.AppendToBuffer("--- Serial error: " + msg.Err.Error() + " ---")
		}
		return m, nil

	case TransferResultMsg:
		handleTransferResult(&m, msg)
		return m, nil

	case tea.KeyMsg:
		newModel, cmd, handled := handleKeyPress(m, msg)
		if handled {
			return newModel, cmd
		}
		if m.CurrentView == state.ViewTerminal && !m.CtrlAPressed {
			var inputCmd tea.Cmd
			m.Input, inputCmd = m.Input.Update(msg)
			return m, inputCmd
		}
		return newModel, cmd
	}

	if m.CurrentView == state.ViewTerminal && !m.CtrlAPressed {
		var cmd tea.Cmd
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// processSerialData processes incoming serial data by splitting it into lines.
// It handles both \n and \r line endings and accumulates incomplete data in the buffer.
func processSerialData(m *model.Model, data string) {
	m.Reader.SerialBuffer += data

	for {
		lineEndIdx := findNextLineEnding(m.Reader.SerialBuffer)
		if lineEndIdx == -1 {
			break
		}

		// Extract the line without the line ending
		line := m.Reader.SerialBuffer[:lineEndIdx]
		m.Reader.SerialBuffer = m.Reader.SerialBuffer[lineEndIdx+1:]

		// Handle CRLF by skipping the following \n if present
		if len(m.Reader.SerialBuffer) > 0 && m.Reader.SerialBuffer[0] == '\n' {
			m.Reader.SerialBuffer = m.Reader.SerialBuffer[1:]
		}

		// Trim any remaining line ending characters and append non-empty lines
		line = strings.Trim(line, "\r\n")
		if line != "" {
			m.AppendToBuffer(line)
		}
	}
}

// findNextLineEnding finds the index of the next line ending (\n or \r).
// Returns -1 if no line ending is found.
func findNextLineEnding(s string) int {
	idxN := strings.Index(s, "\n")
	idxR := strings.Index(s, "\r")

	if idxN == -1 && idxR == -1 {
		return -1
	}

	if idxN == -1 {
		return idxR
	}

	if idxR == -1 {
		return idxN
	}

	// Return whichever comes first
	if idxR < idxN {
		return idxR
	}
	return idxN
}

// handleMouseEvent processes mouse interactions for clickable UI elements.
func handleMouseEvent(m model.Model, msg tea.MouseMsg) (model.Model, tea.Cmd) {
	if m.CurrentView == state.ViewTerminal {
		return handleTerminalMouseEvent(m, msg)
	}

	return handleDialogMouseEvent(m, msg)
}

// handleTerminalMouseEvent handles mouse clicks in the terminal view.
func handleTerminalMouseEvent(m model.Model, msg tea.MouseMsg) (model.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		if m.Zone != nil {
			if zone := m.Zone.Get(model.StatusButtonZoneID); zone != nil && zone.InBounds(msg) {
				if m.Serial.Connected {
					m.Disconnect()
				} else {
					m.Connect()
					if m.Serial.Connected && m.Reader.Reading {
						return m, m.WaitForSerialData()
					}
				}
			}
		}
	}
	return m, nil
}

// handleDialogMouseEvent handles mouse interactions for dragging dialogs.
func handleDialogMouseEvent(m model.Model, msg tea.MouseMsg) (model.Model, tea.Cmd) {
	dialogWidth, dialogHeight := views.GetEstimatedDialogSize(m)
	m.DialogWidth = dialogWidth
	m.DialogHeight = dialogHeight

	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button == tea.MouseButtonLeft {
			dialogX := views.GetDialogX(m)
			dialogY := views.GetDialogY(m)

			// Check if click is in the dialog header (for dragging)
			if msg.X >= dialogX && msg.X < dialogX+dialogWidth &&
				msg.Y >= dialogY && msg.Y < dialogY+DialogHeaderHeight {
				m.IsDragging = true
				m.DragOffsetX = msg.X - dialogX
				m.DragOffsetY = msg.Y - dialogY
			}
		}

	case tea.MouseActionRelease:
		m.IsDragging = false

	case tea.MouseActionMotion:
		if m.IsDragging {
			m.DialogX = clampDialogPosition(msg.X-m.DragOffsetX, 0, m.Width-dialogWidth)
			m.DialogY = clampDialogPosition(msg.Y-m.DragOffsetY, 0, m.Height-dialogHeight)
		}
	}

	return m, nil
}

// clampDialogPosition ensures dialog position stays within screen bounds.
func clampDialogPosition(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// handleKeyPress routes keyboard input to the appropriate handler based on current view.
// Returns the updated model, command, and whether the key was handled.
func handleKeyPress(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd, bool) {
	keyStr := msg.String()

	// Handle Ctrl-A command sequences
	if m.CtrlAPressed {
		return handleCtrlACommands(m, keyStr)
	}

	// Check for Ctrl-A press to enter command mode
	if keyStr == "ctrl+a" {
		m.CtrlAPressed = true
		return m, nil, true
	}

	// Route to view-specific handlers
	switch m.CurrentView {
	case state.ViewTerminal:
		return handleTerminalKey(m, msg)
	case state.ViewMainMenu:
		newModel, cmd := handleMenuKey(m, msg)
		return newModel, cmd, true
	case state.ViewSerialSetup:
		newModel, cmd := handleSerialSetupKey(m, msg)
		return newModel, cmd, true
	case state.ViewHelp:
		newModel, cmd := handleHelpKey(m, msg)
		return newModel, cmd, true
	case state.ViewPortList:
		newModel, cmd := handlePortListKey(m, msg)
		return newModel, cmd, true
	case state.ViewSendFile:
		newModel, cmd := handleSendFileKey(m, msg)
		return newModel, cmd, true
	case state.ViewReceiveFile:
		newModel, cmd := handleReceiveFileKey(m, msg)
		return newModel, cmd, true
	}

	return m, nil, false
}

// handleCtrlACommands processes Ctrl-A command sequences (minicom-style shortcuts).
func handleCtrlACommands(m model.Model, keyStr string) (model.Model, tea.Cmd, bool) {
	m.CtrlAPressed = false

	switch strings.ToLower(keyStr) {
	case "z":
		m.CurrentView = state.ViewHelp
		m.MenuIndex = 0
		m.DialogX = 0
		m.DialogY = 0

	case "c":
		m.TerminalBuffer = make([]string, 0)
		m.CurrentView = state.ViewTerminal

	case "o":
		m.CurrentView = state.ViewMainMenu
		m.MenuIndex = 0
		m.DialogX = 0
		m.DialogY = 0

	case "s":
		m.TransferState.Reset()
		m.TransferState.Direction = transfer.DirectionSend
		m.TransferState.FileInput.Focus()
		m.CurrentView = state.ViewSendFile
		m.DialogX = 0
		m.DialogY = 0

	case "r":
		m.TransferState.Reset()
		m.TransferState.Direction = transfer.DirectionReceive
		m.TransferState.FileInput.Focus()
		m.CurrentView = state.ViewReceiveFile
		m.DialogX = 0
		m.DialogY = 0

	case "l":
		err := m.LoggingState.Toggle("")
		if err != nil {
			m.Serial.LastError = err.Error()
		} else if m.LoggingState.IsEnabled() {
			m.AppendToBuffer("--- Logging started: " + m.LoggingState.GetFilePath() + " ---")
		} else {
			lines, bytes := m.LoggingState.GetStats()
			m.AppendToBuffer(fmt.Sprintf("--- Logging stopped (lines: %d, bytes: %d) ---", lines, bytes))
		}
		m.CurrentView = state.ViewTerminal

	case "q", "x":
		return m, tea.Quit, true

	case "a":
		// Send literal Ctrl-A (ASCII 0x01) to serial port
		if m.Serial.Connected {
			m.WriteSerial([]byte{0x01})
		}
		m.CurrentView = state.ViewTerminal

	default:
		// Unknown command, return to terminal
		m.CurrentView = state.ViewTerminal
	}

	return m, nil, true
}

// handleTerminalKey processes keyboard input in the terminal view.
func handleTerminalKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit, true
	case "enter":
		if m.Input.Value() != "" {
			data := m.Input.Value() + m.GetLineEndingBytes()
			if m.Serial.Connected {
				m.WriteSerial([]byte(data))
			}
			m.AppendToBuffer("> " + m.Input.Value())
			m.Input.SetValue("")
		}
		return m, nil, true
	}
	return m, nil, false
}

// handleMenuKey processes keyboard navigation in the main menu.
func handleMenuKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.MenuIndex > 0 {
			m.MenuIndex--
		}
	case "down", "j":
		if m.MenuIndex < len(m.MenuItems)-1 {
			m.MenuIndex++
		}
	case "enter":
		return selectMenuItem(m)
	case "esc":
		m.CurrentView = state.ViewTerminal
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

// handleSerialSetupKey processes keyboard input in the serial setup dialog.
func handleSerialSetupKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.MenuIndex > 0 {
			m.MenuIndex--
		}
	case "down", "j":
		if m.MenuIndex < menu.SerialSetupMenuItemCount-1 {
			m.MenuIndex++
		}
	case "enter":
		return handleSerialSetupSelect(m)
	case "esc":
		m.CurrentView = state.ViewTerminal
	case "left", "h":
		adjustSerialSetting(&m, -1)
	case "right", "l":
		adjustSerialSetting(&m, 1)
	}
	return m, nil
}

// handleHelpKey processes keyboard input in the help view.
func handleHelpKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", "q", "z":
		m.CurrentView = state.ViewTerminal
	}
	return m, nil
}

// handlePortListKey processes keyboard navigation in the port selection list.
func handlePortListKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.PortIndex > 0 {
			m.PortIndex--
		}
	case "down", "j":
		if m.PortIndex < len(m.AvailablePorts)-1 {
			m.PortIndex++
		}
	case "enter":
		if len(m.AvailablePorts) > 0 {
			m.Serial.Config.Port = m.AvailablePorts[m.PortIndex]
		}
		m.CurrentView = state.ViewSerialSetup
	case "esc":
		m.CurrentView = state.ViewSerialSetup
	}
	return m, nil
}

func handleSendFileKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.TransferState.Reset()
		m.CurrentView = state.ViewTerminal
		return m, nil
	case "tab":
		m.TransferState.FocusNext()
		return m, nil
	case "shift+tab":
		m.TransferState.FocusPrev()
		return m, nil
	case "left", "right":
		if m.TransferState.InputFocus == -1 {
			delta := 1
			if msg.String() == "left" {
				delta = -1
			}
			m.TransferState.CycleProtocol(delta)
		}
		return m, nil
	case "enter":
		return executeSendFile(&m)
	case "ctrl+p":
		m.TransferState.CycleProtocol(1)
		return m, nil
	}

	var cmd tea.Cmd
	m.TransferState.FileInput, cmd = m.TransferState.FileInput.Update(msg)
	return m, cmd
}

func handleReceiveFileKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.TransferState.Reset()
		m.CurrentView = state.ViewTerminal
		return m, nil
	case "tab":
		m.TransferState.FocusNext()
		return m, nil
	case "shift+tab":
		m.TransferState.FocusPrev()
		return m, nil
	case "left", "right":
		if m.TransferState.InputFocus == -1 {
			delta := 1
			if msg.String() == "left" {
				delta = -1
			}
			m.TransferState.CycleProtocol(delta)
		}
		return m, nil
	case "enter":
		return executeReceiveFile(&m)
	case "ctrl+p":
		m.TransferState.CycleProtocol(1)
		return m, nil
	}

	var cmd tea.Cmd
	m.TransferState.FileInput, cmd = m.TransferState.FileInput.Update(msg)
	return m, cmd
}

func executeSendFile(m *model.Model) (model.Model, tea.Cmd) {
	if !m.Serial.Connected || m.Serial.Port == nil {
		m.TransferState.ErrorMsg = "Connect to a serial port first"
		return *m, nil
	}
	if !m.TransferBridge.IsActive() {
		m.TransferState.ErrorMsg = "File transfer bridge unavailable"
		return *m, nil
	}
	if m.TransferBridge.IsTransferring() {
		m.TransferState.ErrorMsg = "Another transfer is already running"
		return *m, nil
	}

	localPath := strings.TrimSpace(m.TransferState.FileInput.Value())
	if localPath == "" {
		m.TransferState.ErrorMsg = "Local file path is required"
		return *m, nil
	}

	absPath, err := transfer.ExpandPath(localPath)
	if err != nil {
		m.TransferState.ErrorMsg = fmt.Sprintf("Invalid path: %v", err)
		return *m, nil
	}
	info, statErr := os.Stat(absPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			m.TransferState.ErrorMsg = "File or directory does not exist"
		} else {
			m.TransferState.ErrorMsg = fmt.Sprintf("Unable to access path: %v", statErr)
		}
		return *m, nil
	}

	resultChan, err := m.TransferBridge.OneTimeUpload([]string{absPath})
	if err != nil {
		m.TransferState.ErrorMsg = fmt.Sprintf("Failed to start transfer: %v", err)
		return *m, nil
	}

	uploadRoot := absPath
	if info != nil && !info.IsDir() {
		uploadRoot = filepath.Dir(absPath)
	}
	if uploadRoot != "" {
		m.TransferBridge.SetDefaultUploadPath(uploadRoot)
	}

	protocolCmd := transfer.GetProtocolCommand(transfer.DirectionSend, m.TransferState.Protocol)
	protoName := m.TransferState.GetProtocolName()
	message := fmt.Sprintf("--- Ready to send %s via %s ---", filepath.Base(absPath), protoName)
	if protocolCmd != "" {
		message += fmt.Sprintf(" Run '%s' on the remote device to begin.", protocolCmd)
	}
	m.AppendToBuffer(message)
	m.BeginTransferWait(transfer.DirectionSend, fmt.Sprintf("Waiting on remote: %s", filepath.Base(absPath)))

	transferPath := absPath
	m.TransferState.Reset()
	m.CurrentView = state.ViewTerminal

	return *m, waitForTransferResult(resultChan, transfer.DirectionSend, transferPath)
}

func executeReceiveFile(m *model.Model) (model.Model, tea.Cmd) {
	if !m.Serial.Connected || m.Serial.Port == nil {
		m.TransferState.ErrorMsg = "Connect to a serial port first"
		return *m, nil
	}
	if !m.TransferBridge.IsActive() {
		m.TransferState.ErrorMsg = "File transfer bridge unavailable"
		return *m, nil
	}
	if m.TransferBridge.IsTransferring() {
		m.TransferState.ErrorMsg = "Another transfer is already running"
		return *m, nil
	}

	localPath := strings.TrimSpace(m.TransferState.FileInput.Value())
	if localPath == "" {
		m.TransferState.ErrorMsg = "Destination directory is required"
		return *m, nil
	}
	absPath, err := transfer.ExpandPath(localPath)
	if err != nil {
		m.TransferState.ErrorMsg = fmt.Sprintf("Invalid path: %v", err)
		return *m, nil
	}
	info, statErr := os.Stat(absPath)
	saveDir := absPath
	if statErr != nil {
		if os.IsNotExist(statErr) {
			if mkErr := os.MkdirAll(absPath, 0o755); mkErr != nil {
				m.TransferState.ErrorMsg = fmt.Sprintf("Failed to create directory: %v", mkErr)
				return *m, nil
			}
		} else {
			m.TransferState.ErrorMsg = fmt.Sprintf("Unable to access path: %v", statErr)
			return *m, nil
		}
	} else if info != nil && !info.IsDir() {
		saveDir = filepath.Dir(absPath)
	}
	if saveDir == "" {
		m.TransferState.ErrorMsg = "Destination directory is required"
		return *m, nil
	}

	m.TransferBridge.SetDefaultDownloadPath(saveDir)
	protoName := m.TransferState.GetProtocolName()
	protocolCmd := transfer.GetProtocolCommand(transfer.DirectionReceive, m.TransferState.Protocol)
	message := fmt.Sprintf("--- Ready to receive files into %s via %s ---", saveDir, protoName)
	if protocolCmd != "" {
		message += fmt.Sprintf(" Ask the remote device to run '%s' with the files you need.", protocolCmd)
	}
	m.AppendToBuffer(message)

	m.TransferState.Reset()
	m.CurrentView = state.ViewTerminal
	return *m, nil
}

func selectMenuItem(m model.Model) (model.Model, tea.Cmd) {
	if m.MenuIndex >= len(m.MenuItems) {
		return m, nil
	}

	item := m.MenuItems[m.MenuIndex]
	switch item.Key {
	case "O":
		m.CurrentView = state.ViewSerialSetup
		m.MenuIndex = 0
		m.DialogX = 0
		m.DialogY = 0
	case "S":
		m.TransferState.Reset()
		m.TransferState.Direction = transfer.DirectionSend
		m.TransferState.FileInput.Focus()
		m.CurrentView = state.ViewSendFile
		m.DialogX = 0
		m.DialogY = 0
	case "R":
		m.TransferState.Reset()
		m.TransferState.Direction = transfer.DirectionReceive
		m.TransferState.FileInput.Focus()
		m.CurrentView = state.ViewReceiveFile
		m.DialogX = 0
		m.DialogY = 0
	case "L":
		err := m.LoggingState.Toggle("")
		if err != nil {
			m.Serial.LastError = err.Error()
		} else if m.LoggingState.IsEnabled() {
			m.AppendToBuffer("--- Logging started: " + m.LoggingState.GetFilePath() + " ---")
		} else {
			lines, bytes := m.LoggingState.GetStats()
			m.AppendToBuffer(fmt.Sprintf("--- Logging stopped (lines: %d, bytes: %d) ---", lines, bytes))
		}
		m.CurrentView = state.ViewTerminal
	case "C":
		m.TerminalBuffer = make([]string, 0)
		m.CurrentView = state.ViewTerminal
	case "Z":
		m.CurrentView = state.ViewHelp
		m.MenuIndex = 0
		m.DialogX = 0
		m.DialogY = 0
	case "X", "Q":
		return m, tea.Quit
	default:
		m.CurrentView = state.ViewTerminal
	}
	return m, nil
}

func handleSerialSetupSelect(m model.Model) (model.Model, tea.Cmd) {
	switch m.MenuIndex {
	case menu.SerialSetupPort:
		m.CurrentView = state.ViewPortList
		m.PortIndex = 0
		m.DialogX = 0
		m.DialogY = 0
		for i, p := range m.AvailablePorts {
			if p == m.Serial.Config.Port {
				m.PortIndex = i
				break
			}
		}
	case menu.SerialSetupConnect:
		if m.Serial.Connected {
			m.Disconnect()
		} else {
			m.Connect()
			if m.Serial.Connected && m.Reader.Reading {
				return m, m.WaitForSerialData()
			}
		}
	}
	return m, nil
}

func adjustSerialSetting(m *model.Model, delta int) {
	baudRates := serial.BaudRates()
	dataBitsOptions := serial.DataBitsOptions()
	stopBitsOptions := serial.StopBitsOptions()
	lineEndingOptions := serial.LineEndingOptions()

	switch m.MenuIndex {
	case menu.SerialSetupBaudRate:
		idx := 0
		for i, b := range baudRates {
			if b == m.Serial.Config.BaudRate {
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
		m.Serial.Config.BaudRate = baudRates[idx]
	case menu.SerialSetupDataBits:
		idx := 0
		for i, d := range dataBitsOptions {
			if d == m.Serial.Config.DataBits {
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
		m.Serial.Config.DataBits = dataBitsOptions[idx]
	case menu.SerialSetupStopBits:
		idx := 0
		for i, s := range stopBitsOptions {
			if s == m.Serial.Config.StopBits {
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
		m.Serial.Config.StopBits = stopBitsOptions[idx]
	case menu.SerialSetupLineEnding:
		idx := 0
		for i, le := range lineEndingOptions {
			if le == m.Serial.Config.LineEnding {
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
		m.Serial.Config.LineEnding = lineEndingOptions[idx]
	}
}

func handleTransferResult(m *model.Model, msg TransferResultMsg) {
	m.ClearTransferWait()
	label := transfer.DirectionLabel(msg.Direction)
	if msg.Err != nil {
		m.AppendToBuffer(fmt.Sprintf("--- %s transfer failed: %v ---", label, msg.Err))
		return
	}
	if msg.Direction == transfer.DirectionSend {
		m.AppendToBuffer(fmt.Sprintf("--- Send complete: %s ---", filepath.Base(msg.Path)))
		return
	}
	m.AppendToBuffer(fmt.Sprintf("--- Receive complete: %s ---", msg.Path))
}

func waitForTransferResult(ch <-chan error, direction transfer.Direction, path string) tea.Cmd {
	return func() tea.Msg {
		if ch == nil {
			return TransferResultMsg{Direction: direction, Path: path, Err: fmt.Errorf("transfer channel unavailable")}
		}
		err, ok := <-ch
		if !ok {
			return TransferResultMsg{Direction: direction, Path: path}
		}
		return TransferResultMsg{Direction: direction, Path: path, Err: err}
	}
}
