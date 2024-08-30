/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var (
	baudrate int
	port     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nanocom",
	Short: "Mighty friendly nano serial communication program",
	Long: `nanocom is a serial communication program.

Used by developers in embedded systems and host machine to communicate serial devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		mode := &serial.Mode{
			BaudRate: baudrate,
		}
		op, err := serial.Open(port, mode)
		if err != nil {
			log.Fatal(err)
		}

		buff := make([]byte, 100)
		for {
			n, err := op.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {
				fmt.Println("\nEOF")
				break
			}
			fmt.Printf("%v", string(buff[:n]))
		}
		// cmd.Help()

	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	rootCmd.Flags().IntVarP(&baudrate, "baudrate", "b", 9200, "set baudrate (ignore the value from config)")

	// if err := rootCmd.MarkFlagRequired("baudrate"); err != nil {
	// 	fmt.Println(err)
	// }
	rootCmd.Flags().StringVarP(&port, "port", "p", "COM2", "set port (ignore the value from config)")

	// list := rootCmd.Flags().StringP("list","l","/dev/tty1","Get available list of serial port")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nanocom.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
