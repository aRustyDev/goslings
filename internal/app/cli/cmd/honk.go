package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Hello returns a greeting for the named person.
func Goodbye(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

var honkCmd = &cobra.Command{
	Use:   "honk",
	Short: "Runs the Goslings tool",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
