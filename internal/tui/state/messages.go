/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package state

import "time"

// TickMsg for periodic updates (like time)
type TickMsg time.Time

// SerialDataMsg for incoming serial data
type SerialDataMsg string

// SerialErrorMsg for serial read errors
type SerialErrorMsg struct {
	Err error
}

// PortsListMsg contains the list of available ports
type PortsListMsg []string

// AutoConnectMsg is sent to trigger auto-connection
type AutoConnectMsg struct{}
