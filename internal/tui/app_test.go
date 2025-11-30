/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp("", 0)

	// App should be properly initialized
	view := app.View()
	if view == "" {
		t.Error("NewApp().View() should return non-empty string")
	}
}

func TestNewAppWithPort(t *testing.T) {
	app := NewApp("COM1", 115200)

	// App should be properly initialized with port
	cmd := app.Init()
	if cmd == nil {
		t.Error("NewApp().Init() should return a command")
	}
}

func TestAppInit(t *testing.T) {
	app := NewApp("", 0)

	cmd := app.Init()
	if cmd == nil {
		t.Error("App.Init() should return a command")
	}
}

func TestAppView(t *testing.T) {
	app := NewApp("", 0)

	view := app.View()
	if view == "" {
		t.Error("App.View() should return non-empty string")
	}
}

func TestAppUpdate(t *testing.T) {
	app := NewApp("", 0)

	// Create a simple message
	type testMsg struct{}

	newModel, cmd := app.Update(testMsg{})

	// Should return a model (App)
	if newModel == nil {
		t.Error("App.Update() should return a model")
	}

	// cmd may or may not be nil depending on the message
	_ = cmd
}
