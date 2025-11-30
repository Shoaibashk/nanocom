/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAdaptiveColors(t *testing.T) {
	colors := []struct {
		name  string
		color lipgloss.AdaptiveColor
	}{
		{"PrimaryColor", PrimaryColor},
		{"SecondaryColor", SecondaryColor},
		{"AccentColor", AccentColor},
		{"ErrorColor", ErrorColor},
		{"SuccessColor", SuccessColor},
		{"WarningColor", WarningColor},
		{"MutedColor", MutedColor},
		{"HighlightColor", HighlightColor},
	}

	for _, c := range colors {
		t.Run(c.name, func(t *testing.T) {
			// Verify both light and dark values are set
			if c.color.Light == "" {
				t.Errorf("%s.Light is empty", c.name)
			}
			if c.color.Dark == "" {
				t.Errorf("%s.Dark is empty", c.name)
			}
		})
	}
}

func TestBaseStyleExists(t *testing.T) {
	// BaseStyle should be defined
	_ = BaseStyle
}

func TestStatusBarStyle(t *testing.T) {
	// StatusBarStyle should have padding
	rendered := StatusBarStyle.Render("test")
	if rendered == "" {
		t.Error("StatusBarStyle should render content")
	}
}

func TestTitleBarStyle(t *testing.T) {
	rendered := TitleBarStyle.Render("Title")
	if rendered == "" {
		t.Error("TitleBarStyle should render content")
	}
}

func TestMenuStyles(t *testing.T) {
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"MenuStyle", MenuStyle},
		{"MenuTitleStyle", MenuTitleStyle},
		{"MenuItemStyle", MenuItemStyle},
		{"MenuSelectedStyle", MenuSelectedStyle},
		{"MenuHotkeyStyle", MenuHotkeyStyle},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			rendered := s.style.Render("test")
			if rendered == "" {
				t.Errorf("%s should render content", s.name)
			}
		})
	}
}

func TestTerminalStyle(t *testing.T) {
	rendered := TerminalStyle.Render("terminal content")
	if rendered == "" {
		t.Error("TerminalStyle should render content")
	}
}

func TestHelpStyle(t *testing.T) {
	rendered := HelpStyle.Render("help text")
	if rendered == "" {
		t.Error("HelpStyle should render content")
	}
}

func TestPanelStyle(t *testing.T) {
	rendered := PanelStyle.Render("panel content")
	if rendered == "" {
		t.Error("PanelStyle should render content")
	}
}

func TestInputStyles(t *testing.T) {
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"InputStyle", InputStyle},
		{"InputPromptStyle", InputPromptStyle},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			rendered := s.style.Render("input")
			if rendered == "" {
				t.Errorf("%s should render content", s.name)
			}
		})
	}
}

func TestLabelStyle(t *testing.T) {
	rendered := LabelStyle.Render("Label:")
	if rendered == "" {
		t.Error("LabelStyle should render content")
	}
}

func TestDialogStyle(t *testing.T) {
	rendered := DialogStyle.Render("dialog content")
	if rendered == "" {
		t.Error("DialogStyle should render content")
	}
}

func TestStatusStyles(t *testing.T) {
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"ErrorStyle", ErrorStyle},
		{"SuccessStyle", SuccessStyle},
		{"WarningStyle", WarningStyle},
		{"ConnectedStyle", ConnectedStyle},
		{"DisconnectedStyle", DisconnectedStyle},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			rendered := s.style.Render("status")
			if rendered == "" {
				t.Errorf("%s should render content", s.name)
			}
		})
	}
}

func TestDividerStyle(t *testing.T) {
	rendered := DividerStyle.Render("---")
	if rendered == "" {
		t.Error("DividerStyle should render content")
	}
}

func TestLogoStyle(t *testing.T) {
	rendered := LogoStyle.Render("nanocom")
	if rendered == "" {
		t.Error("LogoStyle should render content")
	}
}

func TestStylesNotNil(t *testing.T) {
	// Ensure all exported styles can be used
	styles := []lipgloss.Style{
		BaseStyle,
		StatusBarStyle,
		TitleBarStyle,
		MenuStyle,
		MenuTitleStyle,
		MenuItemStyle,
		MenuSelectedStyle,
		MenuHotkeyStyle,
		TerminalStyle,
		HelpStyle,
		PanelStyle,
		InputStyle,
		InputPromptStyle,
		LabelStyle,
		DialogStyle,
		ErrorStyle,
		SuccessStyle,
		WarningStyle,
		ConnectedStyle,
		DisconnectedStyle,
		DividerStyle,
		LogoStyle,
	}

	for i, style := range styles {
		// Just verify each style can render without panic
		result := style.Render("test")
		if result == "" {
			t.Errorf("Style %d rendered empty string", i)
		}
	}
}

func TestStyleWidth(t *testing.T) {
	// Test that width can be set on styles
	styled := StatusBarStyle.Width(80).Render("test")
	if styled == "" {
		t.Error("Setting width should work")
	}
}

func TestStyleCopy(t *testing.T) {
	// Test that styles can be copied and modified
	original := ConnectedStyle
	copied := original.Copy().Background(lipgloss.Color("#000"))

	// Original should be unchanged
	originalRendered := original.Render("test")
	copiedRendered := copied.Render("test")

	// They should be different (copied has background)
	// Just verify both render without error
	if originalRendered == "" || copiedRendered == "" {
		t.Error("Both original and copied styles should render")
	}
}
