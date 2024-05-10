package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRouterInit initializes the routes related to users.
func UserRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewUserRepository(db)

	// Private
	users := router.Group("/users/:companyID")
	{

		// POST endpoint to create a new user
		users.POST("", baseInstance.CreateUser)

		// GET endpoint to retrieve all users for a specific company
		users.GET("", baseInstance.ReadUsers)

		// GET endpoint to retrieve a list of users for a specific company
		users.GET("/list", baseInstance.ReadUsersList)

		// GET endpoint to retrieve the count of users for a specific company
		users.GET("/count", baseInstance.ReadUsersCount)

		// GET endpoint to retrieve details of a specific user
		users.GET("/:ID", baseInstance.ReadUser)

		// PUT endpoint to update details of a specific user
		users.PUT("/:ID", baseInstance.UpdateUser)

		// DELETE endpoint to delete a specific user
		users.DELETE("/:ID", baseInstance.DeleteUser)
		

		users.GET("/gender", baseInstance.GenderPercentage)
	}


}
