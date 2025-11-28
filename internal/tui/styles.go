/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for minicom-like appearance
var (
	// Primary colors
	primaryColor   = lipgloss.Color("#00AAFF")
	secondaryColor = lipgloss.Color("#FFFFFF")
	bgColor        = lipgloss.Color("#000033")
	accentColor    = lipgloss.Color("#FFFF00")
	errorColor     = lipgloss.Color("#FF0000")
	successColor   = lipgloss.Color("#00FF00")
	menuBgColor    = lipgloss.Color("#000066")
	menuBorderCol  = lipgloss.Color("#0088FF")
	highlightColor = lipgloss.Color("#00FFFF")
)

// Styles for various UI components
var (
	// Base styles
	baseStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(secondaryColor)

	// Status bar at the bottom
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0044AA")).
			Foreground(secondaryColor).
			Bold(true).
			Padding(0, 1)

	// Title bar at the top
	titleBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0044AA")).
			Foreground(accentColor).
			Bold(true).
			Padding(0, 1)

	// Menu styles
	menuStyle = lipgloss.NewStyle().
			Background(menuBgColor).
			Foreground(secondaryColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(menuBorderCol).
			Padding(1, 2)

	menuTitleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Underline(true).
			MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			PaddingLeft(1)

	menuSelectedStyle = lipgloss.NewStyle().
				Background(highlightColor).
				Foreground(lipgloss.Color("#000000")).
				Bold(true).
				PaddingLeft(1)

	menuHotkeyStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// Terminal view style
	terminalStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(successColor)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)

	// Border style for panels
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(menuBorderCol).
			Padding(0, 1)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#003366")).
			Foreground(secondaryColor).
			Padding(0, 1)

	// Label style
	labelStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Dialog box style
	dialogStyle = lipgloss.NewStyle().
			Background(menuBgColor).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(accentColor).
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
)
