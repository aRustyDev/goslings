package cmd

import (
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "middle command for dumping data from endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
