// Goodbye returns a farewell for the named person.
package cmd

import "fmt"

func Goodbye(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Goodbye, %v. Siyonara!", name)
	return message
}
