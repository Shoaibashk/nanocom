/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/shoaibashk/nanocom/internal/tui/model"
	"github.com/shoaibashk/nanocom/internal/tui/state"
	"github.com/shoaibashk/nanocom/internal/tui/styles"
)

const (
	// UI layout constants
	titleBarHeight  = 1
	statusBarHeight = 1
	inputLineHeight = 1
	borderHeight    = 1
	totalUIOverhead = titleBarHeight + statusBarHeight + inputLineHeight + borderHeight

	// Rendering constants
	emptyPortDisplay = "(none selected)"
)

// View renders the appropriate view based on the current application state.
// It wraps the view with bubble zone scanning for clickable regions.
func View(m model.Model) string {
	if m.Width == 0 || m.Height == 0 {
		return "Initializing..."
	}

	switch m.CurrentView {
	case state.ViewTerminal:
		return renderWithZones(m, RenderTerminalView(m))
	case state.ViewMainMenu:
		return renderWithZones(m, RenderMainMenu(m))
	case state.ViewSerialSetup:
		return renderWithZones(m, RenderSerialSetup(m))
	case state.ViewHelp:
		return renderWithZones(m, RenderHelpMenu(m))
	case state.ViewPortList:
		return renderWithZones(m, RenderPortList(m))
	case state.ViewSendFile:
		return renderWithZones(m, RenderSendFileView(m))
	case state.ViewReceiveFile:
		return renderWithZones(m, RenderReceiveFileView(m))
	default:
		return renderWithZones(m, RenderTerminalView(m))
	}
}

// renderWithZones scans the view for clickable zones if zone manager is available.
func renderWithZones(m model.Model, view string) string {
	if m.Zone != nil {
		return m.Zone.Scan(view)
	}
	return view
}

// RenderTerminalView renders the main terminal view with title bar, terminal content, input, and status bar.
func RenderTerminalView(m model.Model) string {
	// Title bar with branding
	logoText := "◆ nanocom"
	versionText := "v1.0"
	title := styles.LogoStyle.Render(logoText) + " " + styles.HelpStyle.Render(versionText)
	titleBar := styles.TitleBarStyle.Width(m.Width).Render(title)

	// Calculate available height for terminal content
	termHeight := m.Height - totalUIOverhead - 3 // Extra space for safety

	// Render terminal content with line wrapping
	termContent := renderTerminalContent(m, termHeight)
	terminal := styles.TerminalStyle.Width(m.Width).Height(termHeight).Render(termContent)

	// Input line with prompt
	prompt := styles.InputPromptStyle.Render("❯ ")
	inputField := m.Input.View()
	inputLine := styles.InputStyle.Width(m.Width - 4).Render(prompt + inputField)

	// Status bar
	status := RenderStatusBar(m)

	return lipgloss.JoinVertical(lipgloss.Left, titleBar, terminal, inputLine, status)
}

// renderTerminalContent renders the terminal buffer content with styling.
func renderTerminalContent(m model.Model, termHeight int) string {
	var content strings.Builder

	// Calculate visible lines
	startIdx := 0
	if len(m.TerminalBuffer) > termHeight {
		startIdx = len(m.TerminalBuffer) - termHeight
	}

	// Render each line with appropriate styling
	for i := startIdx; i < len(m.TerminalBuffer); i++ {
		line := m.TerminalBuffer[i]

		// Style different types of messages
		if strings.HasPrefix(line, "> ") {
			content.WriteString(styles.InputPromptStyle.Render(">") + " " + line[2:])
		} else if strings.HasPrefix(line, "---") {
			content.WriteString(styles.HelpStyle.Render(line))
		} else {
			content.WriteString(line)
		}

		if i < len(m.TerminalBuffer)-1 {
			content.WriteString("\n")
		}
	}

	// Pad with empty lines to fill terminal height
	lineCount := len(m.TerminalBuffer) - startIdx
	for i := lineCount; i < termHeight; i++ {
		content.WriteString("\n")
	}

	return content.String()
}

// RenderStatusBar renders the bottom status bar with port info, connection status, and help hints.
func RenderStatusBar(m model.Model) string {
	portInfo := formatPortInfo(m)
	connBadge := formatConnectionBadge(m)
	logBadge := formatLoggingBadge(m)
	transferBadge := formatTransferBadge(m)
	timeStr := time.Now().Format("15:04:05")
	hint := styles.HelpStyle.Render("Ctrl+A Z for help")

	// Build left and right sections
	left := portInfo + "  " + connBadge + logBadge + transferBadge
	right := hint + "  " + styles.HelpStyle.Render(timeStr)

	// Create clickable connect/disconnect button
	button := formatConnectionButton(m)
	buttonRendered := button
	if m.Zone != nil {
		buttonRendered = m.Zone.Mark(model.StatusButtonZoneID, button)
	}

	// Calculate spacing
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	buttonWidth := lipgloss.Width(button)
	totalContent := leftWidth + buttonWidth + rightWidth
	gap := m.Width - totalContent - 6
	if gap < 2 {
		gap = 2
	}
	leftGap := gap / 2
	rightGap := gap - leftGap

	statusContent := left + strings.Repeat(" ", leftGap) + buttonRendered + strings.Repeat(" ", rightGap) + right
	return styles.StatusBarStyle.Width(m.Width).Render(statusContent)
}

// formatPortInfo returns formatted serial port configuration string.
func formatPortInfo(m model.Model) string {
	parityChar := "N"
	if len(m.Serial.Config.Parity) > 0 {
		parityChar = string(m.Serial.Config.Parity[0])
	}

	portDisplay := m.Serial.Config.Port
	if portDisplay == "" {
		portDisplay = "No port"
	}

	return fmt.Sprintf("⚡ %s • %d %d%s%d",
		portDisplay, m.Serial.Config.BaudRate, m.Serial.Config.DataBits, parityChar, m.Serial.Config.StopBits)
}

// formatConnectionBadge returns styled connection status badge.
func formatConnectionBadge(m model.Model) string {
	if m.Serial.Connected {
		return styles.ConnectedStyle.Render("● ONLINE")
	}
	return styles.DisconnectedStyle.Render("○ OFFLINE")
}

// formatLoggingBadge returns styled logging indicator if logging is active.
func formatLoggingBadge(m model.Model) string {
	if m.LoggingState.IsEnabled() {
		return styles.WarningStyle.Render(" 📝 LOG")
	}
	return ""
}

// formatTransferBadge returns styled transfer status if transfer is pending.
func formatTransferBadge(m model.Model) string {
	if status := m.CurrentTransferStatus(); status != "" {
		return styles.WarningStyle.Render(" " + status)
	}
	return ""
}

// formatConnectionButton returns styled connect/disconnect button.
func formatConnectionButton(m model.Model) string {
	if m.Serial.Connected {
		return styles.DisconnectedStyle.Copy().Background(lipgloss.Color("#333")).Render("[ Disconnect ]")
	}
	return styles.ConnectedStyle.Copy().Background(lipgloss.Color("#333")).Render("[ Connect ]")
}

// RenderMainMenu renders the main menu
func RenderMainMenu(m model.Model) string {
	bg := RenderTerminalView(m)

	var menuContent strings.Builder
	menuContent.WriteString(styles.MenuTitleStyle.Render("◆ nanocom Command Summary"))
	menuContent.WriteString("\n\n")

	for i, item := range m.MenuItems {
		if i == m.MenuIndex {
			line := styles.MenuSelectedStyle.Render(fmt.Sprintf("[%s] %s", item.Key, item.Label))
			menuContent.WriteString(line)
		} else {
			hotkey := styles.MenuHotkeyStyle.Render(fmt.Sprintf("[%s]", item.Key))
			line := styles.MenuItemStyle.Render(fmt.Sprintf("%s %s", hotkey, item.Label))
			menuContent.WriteString(line)
		}
		menuContent.WriteString("\n")
	}

	menuContent.WriteString("\n")
	menuContent.WriteString(styles.HelpStyle.Render("↑/↓ navigate • enter select • esc close"))

	menuBox := styles.MenuStyle.Render(menuContent.String())
	return centerDialog(m, bg, menuBox)
}

// RenderSerialSetup renders the serial port configuration dialog.
func RenderSerialSetup(m model.Model) string {
	bg := RenderTerminalView(m)

	var content strings.Builder
	content.WriteString(styles.MenuTitleStyle.Render("Serial Port Settings"))
	content.WriteString("\n\n")

	portDisplay := m.Serial.Config.Port
	if portDisplay == "" {
		portDisplay = emptyPortDisplay
	}

	settings := []struct {
		label string
		value string
	}{
		{"Serial Device", portDisplay},
		{"Baud Rate", fmt.Sprintf("%d", m.Serial.Config.BaudRate)},
		{"Data Bits", fmt.Sprintf("%d", m.Serial.Config.DataBits)},
		{"Stop Bits", fmt.Sprintf("%d", m.Serial.Config.StopBits)},
		{"Line Ending", m.Serial.GetLineEndingDisplay()},
		{m.GetConnectLabel(), ""},
	}

	// Render settings with alignment
	maxLabelLen := findMaxLabelLength(settings)
	for i, s := range settings {
		renderSetting(&content, s.label, s.value, maxLabelLen, i == m.MenuIndex)
	}

	content.WriteString("\n")
	content.WriteString(styles.HelpStyle.Render("↑/↓ navigate • ←/→ change • enter select • esc back"))

	if m.Serial.LastError != "" {
		content.WriteString("\n\n")
		content.WriteString(styles.ErrorStyle.Render("⚠ " + m.Serial.LastError))
	}

	menuBox := styles.MenuStyle.Render(content.String())
	return centerDialog(m, bg, menuBox)
}

// findMaxLabelLength finds the longest label for alignment.
func findMaxLabelLength(settings []struct{ label, value string }) int {
	maxLen := 0
	for _, s := range settings {
		if len(s.label) > maxLen {
			maxLen = len(s.label)
		}
	}
	return maxLen
}

// renderSetting renders a single setting line with proper styling.
func renderSetting(sb *strings.Builder, label, value string, maxLabelLen int, selected bool) {
	paddedLabel := fmt.Sprintf("%-*s", maxLabelLen, label)

	if selected {
		if value != "" {
			sb.WriteString(styles.MenuSelectedStyle.Render(fmt.Sprintf("%s : %s", paddedLabel, value)))
		} else {
			sb.WriteString(styles.MenuSelectedStyle.Render(paddedLabel))
		}
	} else {
		styledLabel := styles.LabelStyle.Render(paddedLabel)
		if value != "" {
			sb.WriteString(styles.MenuItemStyle.Render(fmt.Sprintf("%s : %s", styledLabel, value)))
		} else {
			sb.WriteString(styles.MenuItemStyle.Render(styledLabel))
		}
	}
	sb.WriteString("\n")
}

// RenderHelpMenu renders the help menu
func RenderHelpMenu(m model.Model) string {
	bg := RenderTerminalView(m)

	var content strings.Builder
	content.WriteString(styles.MenuTitleStyle.Render("◆ nanocom Help"))
	content.WriteString("\n\n")

	for _, item := range m.HelpMenuItems {
		hotkey := styles.MenuHotkeyStyle.Render(fmt.Sprintf("%-10s", item.Key))
		label := item.Label
		content.WriteString(fmt.Sprintf("  %s  %s\n", hotkey, label))
	}

	content.WriteString("\n")
	content.WriteString(styles.HelpStyle.Render("Press any key to close"))

	menuBox := styles.MenuStyle.Render(content.String())
	return centerDialog(m, bg, menuBox)
}

// RenderPortList renders the port list dialog
func RenderPortList(m model.Model) string {
	bg := RenderTerminalView(m)

	var content strings.Builder
	content.WriteString(styles.MenuTitleStyle.Render("Select Serial Port"))
	content.WriteString("\n\n")

	if len(m.AvailablePorts) == 0 {
		content.WriteString(styles.ErrorStyle.Render("⚠ No serial ports found"))
	} else {
		for i, port := range m.AvailablePorts {
			if i == m.PortIndex {
				content.WriteString(styles.MenuSelectedStyle.Render(" ● " + port + " "))
			} else {
				content.WriteString(styles.MenuItemStyle.Render("   " + port))
			}
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(styles.HelpStyle.Render("↑/↓ navigate • enter select • esc back"))

	menuBox := styles.MenuStyle.Render(content.String())
	return centerDialog(m, bg, menuBox)
}

// centerDialog positions and overlays a dialog box, centering it if no position is set.
func centerDialog(m model.Model, bg, menuBox string) string {
	menuWidth := lipgloss.Width(menuBox)
	menuHeight := lipgloss.Height(menuBox)

	x := m.DialogX
	y := m.DialogY
	if x == 0 && y == 0 {
		x = (m.Width - menuWidth) / 2
		y = (m.Height - menuHeight) / 2
	}

	return OverlayAtWithSize(m, bg, menuBox, x, y, menuWidth, menuHeight)
}

// RenderHelpMenu renders the help menu dialog.
func RenderSendFileView(m model.Model) string {
	return renderTransferDialog(m, "Send Files", "enter: send | esc: cancel")
}

// RenderReceiveFileView renders the receive file dialog
func RenderReceiveFileView(m model.Model) string {
	return renderTransferDialog(m, "Receive Files", "enter: receive | esc: cancel")
}

// renderTransferDialog renders the file transfer dialog (send or receive).
func renderTransferDialog(m model.Model, title, helpHint string) string {
	bg := RenderTerminalView(m)

	var content strings.Builder
	content.WriteString(styles.MenuTitleStyle.Render(title))
	content.WriteString("\n\n")

	protocolLabel := styles.LabelStyle.Render("Protocol:")
	protocolValue := m.TransferState.GetProtocolName()
	content.WriteString(fmt.Sprintf("  %s %s (Ctrl+P to change)\n\n", protocolLabel, protocolValue))

	label := styles.LabelStyle.Render("Local File:")
	content.WriteString(fmt.Sprintf("  %s %s\n", styles.MenuSelectedStyle.Render(">"), label))
	content.WriteString(fmt.Sprintf("    %s\n\n", m.TransferState.FileInput.View()))

	if m.TransferState.ErrorMsg != "" {
		content.WriteString(styles.ErrorStyle.Render(m.TransferState.ErrorMsg))
		content.WriteString("\n")
	}
	if m.TransferState.StatusMsg != "" {
		content.WriteString(styles.SuccessStyle.Render(m.TransferState.StatusMsg))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(styles.HelpStyle.Render(helpHint))

	menuBox := styles.MenuStyle.Render(content.String())
	return centerDialog(m, bg, menuBox)
}

// OverlayAtWithSize overlays content at a specific position on the background.
// It properly handles ANSI escape sequences and maintains text styling.
func OverlayAtWithSize(m model.Model, bg, overlay string, x, y, width, height int) string {
	bgLines := strings.Split(bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for len(bgLines) < m.Height {
		bgLines = append(bgLines, strings.Repeat(" ", m.Width))
	}

	for i, overlayLine := range overlayLines {
		lineY := y + i
		if lineY >= 0 && lineY < len(bgLines) {
			bgLine := bgLines[lineY]

			for len(bgLine) < x+lipgloss.Width(overlayLine) {
				bgLine += " "
			}

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

// GetEstimatedDialogSize returns the estimated width and height for a dialog.
// Used for positioning and dragging calculations.
func GetEstimatedDialogSize(m model.Model) (int, int) {
	switch m.CurrentView {
	case state.ViewMainMenu:
		return 40, len(m.MenuItems) + 8
	case state.ViewSerialSetup:
		return 50, 12
	case state.ViewHelp:
		return 45, len(m.HelpMenuItems) + 8
	case state.ViewPortList:
		ports := len(m.AvailablePorts)
		if ports == 0 {
			ports = 1
		}
		return 40, ports + 8
	case state.ViewSendFile, state.ViewReceiveFile:
		return 60, 14
	default:
		return 40, 15
	}
}

// GetDialogX returns the dialog's X position, centering it if not explicitly set.
func GetDialogX(m model.Model) int {
	if m.DialogX == 0 {
		width, _ := GetEstimatedDialogSize(m)
		return (m.Width - width) / 2
	}
	return m.DialogX
}

// GetDialogY returns the dialog's Y position, centering it if not explicitly set.
func GetDialogY(m model.Model) int {
	if m.DialogY == 0 {
		_, height := GetEstimatedDialogSize(m)
		return (m.Height - height) / 2
	}
	return m.DialogY
}
