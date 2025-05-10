package main

import (
	"fmt"
	"os"

	"github.com/arustydev/goslings/internal/app/cli/cmd"
	"github.com/arustydev/goslings/internal/app/tui"
	log "github.com/sirupsen/logrus"
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
	log.Println("Hello logrus!")
	fmt.Println(Hello("name string"))
}

func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	cmd.Goodbye(name)
	return message
}
