/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/shoaibashk/nanocom/internal/tui/model"
	"github.com/shoaibashk/nanocom/internal/tui/update"
	"github.com/shoaibashk/nanocom/internal/tui/views"
)

// App wraps the model to implement tea.Model interface
type App struct {
	model model.Model
}

// NewApp creates a new TUI application
func NewApp(port string, baudRate int) App {
	return App{
		model: model.New(port, baudRate),
	}
}

// Init initializes the application
func (a App) Init() tea.Cmd {
	return a.model.Init()
}

// Update handles messages
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := update.Update(a.model, msg)
	a.model = newModel
	return a, cmd
}

// View renders the application
func (a App) View() string {
	return views.View(a.model)
}

// Run starts the TUI
func Run(port string, baudRate int) error {
	app := NewApp(port, baudRate)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
