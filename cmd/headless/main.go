package main

import (
	"fmt"
	"os"

	"goslings/internal/app/tui"
	"goslings/internal/auth"
)

func main() {
	p := tui.NewTui()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Printf(Hello("name string"))
}

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	auth.Goodbye(name)
	return message
}
