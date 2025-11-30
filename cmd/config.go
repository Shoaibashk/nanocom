/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com

Note: This command is currently unimplemented and reserved for future configuration management.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd manages nanocom configuration (unimplemented)
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage nanocom configuration",
	Long:  `View and modify nanocom configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Configuration management is not yet implemented.")
		fmt.Println("Use the TUI interface (Ctrl-A O) to configure serial settings.")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
