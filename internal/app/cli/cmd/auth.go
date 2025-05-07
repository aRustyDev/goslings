package cmd

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "middle command for authenticating to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
