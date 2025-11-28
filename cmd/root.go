/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/shoaibashk/nanocom/internal/tui"
	"github.com/spf13/cobra"
)

var (
	baudrate int
	port     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nanocom",
	Short: "Mighty friendly nano serial communication program",
	Long: `nanocom is a cross-platform serial communication program with a TUI interface.

It is designed to resemble minicom and provides an intuitive terminal interface
for communicating with serial devices such as embedded systems, modems, and routers.

Usage:
  nanocom                    Start with TUI (use Ctrl-A Z for help)
  nanocom -p /dev/ttyUSB0    Start with specific port
  nanocom -b 115200          Start with specific baud rate`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(port, baudrate); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().IntVarP(&baudrate, "baudrate", "b", 9600, "set baudrate")
	rootCmd.Flags().StringVarP(&port, "port", "p", "", "set port (e.g., /dev/ttyUSB0 or COM1)")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nanocom.yaml)")
}
