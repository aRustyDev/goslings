package main

import (
	"goslings/internal/app/cli/cmd"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// https://masteringbackend.com/posts/gin-framework#getting-started-with-gin
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	// Create a new Gin router
	router := gin.Default()

	// Define a route for the root URL
	router.GET("/", func(c *gin.Context) {
		c.String(200, cmd.Goodbye("name"))
	})

	// Run the server on port 8080
	router.Run(":8080")
}
