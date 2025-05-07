package cmd

import (
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "convenience command for working with Gosling configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
