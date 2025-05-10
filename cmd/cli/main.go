package main

import (
	"context"
	"fmt"
	"time"

	"github.com/arustydev/goslings/internal/app/cli/cmd"
)

func main() {
	cmd.Goodbye("from the cli!")
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd.Execute(ctx)
}

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	cmd.Goodbye(name)
	return message
}
