package main

import (
	"fmt"
	"os"

	"github.com/devgymbr/fit/command"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fit",
	Short: "fit is a minimalistc git cli tool",
	Long:  "fit is a minimalistc git cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func main() {
	rootCmd.AddCommand(command.Init)
	rootCmd.AddCommand(command.Add)
	rootCmd.AddCommand(command.Status)
	rootCmd.AddCommand(command.Commit)
	rootCmd.AddCommand(command.Log)
	rootCmd.AddCommand(command.Show)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
