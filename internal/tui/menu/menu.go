/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package menu

// Item represents a menu option
type Item struct {
	Key   string
	Label string
}

// Serial setup menu item indices
const (
	SerialSetupPort = iota
	SerialSetupBaudRate
	SerialSetupDataBits
	SerialSetupStopBits
	SerialSetupLineEnding
	SerialSetupConnect
	SerialSetupMenuItemCount
)

// DefaultMainMenuItems returns the default main menu items
func DefaultMainMenuItems() []Item {
	return []Item{
		{Key: "O", Label: "Serial pOrt Settings"},
		{Key: "S", Label: "Send files"},
		{Key: "R", Label: "Receive files"},
		{Key: "C", Label: "Clear Screen"},
		{Key: "L", Label: "capture to fiLe (Log)"},
		{Key: "Z", Label: "Help (show this menu)"},
		{Key: "X", Label: "eXit and reset"},
		{Key: "Q", Label: "Quit with no reset"},
	}
}

// DefaultHelpMenuItems returns the default help menu items
func DefaultHelpMenuItems() []Item {
	return []Item{
		{Key: "^A Z", Label: "Show Help Menu"},
		{Key: "^A O", Label: "Serial Port Settings"},
		{Key: "^A C", Label: "Clear Screen"},
		{Key: "^A S", Label: "Send Files"},
		{Key: "^A R", Label: "Receive Files"},
		{Key: "^A L", Label: "Toggle Logging"},
		{Key: "^A Q", Label: "Quit"},
		{Key: "^A X", Label: "Exit"},
	}
}
