/*

	Package utils provides utility functions and configurations for the application.

	The InitRouter function initializes and returns a Gin router with specific configurations.

	Functions:
	- InitRouter(): Initializes and returns a Gin router with recovery and logging.

	Dependencies:
	- "github.com/gin-gonic/gin": Gin framework for building web applications in Go.

	Usage:
	- Call InitRouter() to obtain a configured Gin router for the application.

	Note:
	- The function sets the Gin framework mode to "DebugMode".
	- The router includes recovery middleware for handling panics, logging middleware for request logging.

	Last update :
	01/02/2024 10:22

*/

package utils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes and returns a Gin router with specific configurations.
func InitRouter() *gin.Engine {

	gin.SetMode(gin.DebugMode)

	router := gin.New()

	router.Use(gin.Recovery(), gin.Logger())

	// Cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	return router
}
