package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthRouterInit initializes the routes related to authentication.
func AuthRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Public
	auth := router.Group("/auth")
	{

		// POST endpoint for user sign-in
		auth.POST("/signin", baseInstance.SigninUser)

		// POST endpoint for user sign-up
		auth.POST("/signup", baseInstance.SignupUser)
	}
}
