package cmd

import (
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "custom Help command for Goslings",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
