package main

import (
	"fmt"

	"goslings/internal/app/cli/cmd"
	"goslings/internal/auth"
)

func main() {
	cmd.Goodbye("from the cli!")
}

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	auth.Goodbye(name)
	return message
}
