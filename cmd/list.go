/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

// listCmd lists all available serial ports
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available serial ports",
	Long:  `Display all available serial ports on the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		ports, err := serial.GetPortsList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing ports: %v\n", err)
			os.Exit(1)
		}

		if len(ports) == 0 {
			fmt.Println("No serial ports found.")
			return
		}

		fmt.Printf("Available serial ports (%d):\n", len(ports))
		for _, port := range ports {
			fmt.Printf("  - %s\n", port)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
