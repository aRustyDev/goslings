package main

import (
	"goslings/internal/auth"

	"github.com/gin-gonic/gin"
)

// https://masteringbackend.com/posts/gin-framework#getting-started-with-gin
func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route for the root URL
	router.GET("/", func(c *gin.Context) {
		c.String(200, auth.Goodbye("name"))
	})

	// Run the server on port 8080
	router.Run(":8080")
}
