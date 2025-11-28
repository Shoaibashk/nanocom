/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Stock Charm color palette using adaptive colors
var (
	primaryColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	secondaryColor = lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C1C1"}
	accentColor    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	errorColor     = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF6666"}
	successColor   = lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#00FF00"}
	warningColor   = lipgloss.AdaptiveColor{Light: "#FFAA00", Dark: "#FFCC00"}
	mutedColor     = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#626262"}
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)

// Styles for various UI components - Stock Charm styling
var (
	// Base styles
	baseStyle = lipgloss.NewStyle()

	// Status bar at the bottom
	statusBarStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Title bar at the top
	titleBarStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	// Menu styles
	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	menuTitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	menuSelectedStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				PaddingLeft(2)

	menuHotkeyStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Bold(true)

	// Terminal view style
	terminalStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Border style for panels
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	// Input prompt style
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// Label style
	labelStyle = lipgloss.NewStyle().
			Bold(true)

	// Dialog box style
	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Align(lipgloss.Center)

	// Error text style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Success text style
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Warning text style
	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	// Connected status badge
	connectedStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)

	// Disconnected status badge
	disconnectedStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true).
				Padding(0, 1)

	// Divider style
	dividerStyle = lipgloss.NewStyle()

	// Logo/brand style
	logoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)
)
