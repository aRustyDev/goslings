package main

import (
	"fmt"
	"goslings/internal/app/tui"
	"goslings/internal/auth"
	"os"
)

// https://www.taranveerbains.ca/blog/13-making-a-tui-with-go
// https://github.com/charmbracelet/bubbletea
// https://github.com/charmbracelet/lipgloss
// https://github.com/charmbracelet/bubbles
func main() {
	p := tui.NewTui()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Printf(Hello("name string"))
}

func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	auth.Goodbye(name)
	return message
}
