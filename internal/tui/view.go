package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI based on the current view state.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	switch m.currentView {
	case ViewTerminal:
		return m.renderWithZones(m.renderTerminalView())
	case ViewMainMenu:
		return m.renderWithZones(m.renderMainMenu())
	case ViewSerialSetup:
		return m.renderWithZones(m.renderSerialSetup())
	case ViewHelp:
		return m.renderWithZones(m.renderHelpMenu())
	case ViewPortList:
		return m.renderWithZones(m.renderPortList())
	case ViewSendFile:
		return m.renderWithZones(m.renderSendFileView())
	case ViewReceiveFile:
		return m.renderWithZones(m.renderReceiveFileView())
	default:
		return m.renderWithZones(m.renderTerminalView())
	}
}

func (m Model) renderWithZones(view string) string {
	if m.zone != nil {
		return m.zone.Scan(view)
	}
	return view
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

	transferBadge := ""
	if status := m.currentTransferStatus(); status != "" {
		transferBadge = warningStyle.Render(" " + status)
	}

	timeStr := time.Now().Format("15:04:05")
	hint := helpStyle.Render("Ctrl+A Z for help")

	// Build status bar with proper spacing
	left := portInfo + "  " + connBadge + logBadge + transferBadge
	right := hint + "  " + helpStyle.Render(timeStr)

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)

	// Create clickable button
	var button string
	if m.connected {
		button = disconnectedStyle.Copy().Background(lipgloss.Color("#333")).Render("[ Disconnect ]")
	} else {
		button = connectedStyle.Copy().Background(lipgloss.Color("#333")).Render("[ Connect ]")
	}
	buttonWidth := lipgloss.Width(button)
	buttonRendered := button
	if m.zone != nil {
		buttonRendered = m.zone.Mark(statusButtonZoneID, button)
	}

	totalContent := leftWidth + buttonWidth + rightWidth
	gap := m.width - totalContent - 6
	if gap < 2 {
		gap = 2
	}
	leftGap := gap / 2
	rightGap := gap - leftGap

	statusContent := left + strings.Repeat(" ", leftGap) + buttonRendered + strings.Repeat(" ", rightGap) + right
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
	return m.renderTransferDialog("Send Files", "enter: send | esc: cancel")
}

func (m Model) renderReceiveFileView() string {
	return m.renderTransferDialog("Receive Files", "enter: receive | esc: cancel")
}

// renderTransferDialog draws the shared transfer dialog for send/receive flows.
func (m Model) renderTransferDialog(title, helpHint string) string {
	bg := m.renderTerminalView()

	var content strings.Builder
	content.WriteString(menuTitleStyle.Render(title))
	content.WriteString("\n\n")

	// Protocol selector
	protocolLabel := labelStyle.Render("Protocol:")
	protocolValue := m.transferState.GetProtocolName()
	content.WriteString(fmt.Sprintf("  %s %s (Ctrl+P to change)\n\n", protocolLabel, protocolValue))

	// Input field - only local file
	label := labelStyle.Render("Local File:")
	content.WriteString(fmt.Sprintf("  %s %s\n", menuSelectedStyle.Render(">"), label))
	content.WriteString(fmt.Sprintf("    %s\n\n", m.transferState.fileInput.View()))

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
	content.WriteString(helpStyle.Render(helpHint))

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
