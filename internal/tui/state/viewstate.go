/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package state

// ViewState represents the current view/mode of the TUI
type ViewState int

const (
	ViewTerminal ViewState = iota
	ViewMainMenu
	ViewSerialSetup
	ViewHelp
	ViewPortList
	ViewSendFile
	ViewReceiveFile
)
