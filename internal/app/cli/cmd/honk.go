package cmd

import (
	"github.com/spf13/cobra"
)

var honkCmd = &cobra.Command{
	Use:   "honk",
	Short: "Runs the Goslings tool",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
