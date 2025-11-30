/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package serial

import "go.bug.st/serial"

// GetPortsList returns the list of available serial ports
func GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}
