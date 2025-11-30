/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Stock Charm color palette using adaptive colors
var (
	PrimaryColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	SecondaryColor = lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C1C1"}
	AccentColor    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	ErrorColor     = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF6666"}
	SuccessColor   = lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#00FF00"}
	WarningColor   = lipgloss.AdaptiveColor{Light: "#FFAA00", Dark: "#FFCC00"}
	MutedColor     = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#626262"}
	HighlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)

// Styles for various UI components - Stock Charm styling
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle()

	// Status bar at the bottom
	StatusBarStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Title bar at the top
	TitleBarStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	// Menu styles
	MenuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	MenuTitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	MenuSelectedStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				PaddingLeft(2)

	MenuHotkeyStyle = lipgloss.NewStyle().
			Foreground(HighlightColor).
			Bold(true)

	// Terminal view style
	TerminalStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// Border style for panels
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	// Input field style
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	// Input prompt style
	InputPromptStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true)

	// Label style
	LabelStyle = lipgloss.NewStyle().
			Bold(true)

	// Dialog box style
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Align(lipgloss.Center)

	// Error text style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	// Success text style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	// Warning text style
	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	// Connected status badge
	ConnectedStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Padding(0, 1)

	// Disconnected status badge
	DisconnectedStyle = lipgloss.NewStyle().
				Foreground(ErrorColor).
				Bold(true).
				Padding(0, 1)

	// Divider style
	DividerStyle = lipgloss.NewStyle()

	// Logo/brand style
	LogoStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)
)
