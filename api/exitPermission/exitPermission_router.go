package exitPermission

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExitPermissionRouterInit initializes the routes related to exitPermission.
func ExitPermissionRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewExitPermissionRepository(db)

	// Private
	exitPermission := router.Group("/exitPermission")
	{

		//Create a new exitPermission demande
		exitPermission.POST("", baseInstance.AddExitPermission)

		// GET endpoint to retrieve all exitPermission for a specific user
		exitPermission.GET("/:userID", baseInstance.ReadAllExitPermission)

		// GET endpoint to retrieve the count of exitPermission for a specific user
		exitPermission.GET("/user/:userID/count", baseInstance.ReadExitPermissionCount)

		// GET endpoint to retrieve details of a specific exitPermission for a specific user
		exitPermission.GET("/user/:userID/:ID", baseInstance.ReadOneExitPermission)

		// PUT endpoint to update the details of a specific exitPermission for a specific user
		exitPermission.PUT("/user/:userID/:ID", baseInstance.UpdateExitPermission)

		// DELETE endpoint to delete a specific exitPermission for a specific user
		exitPermission.DELETE("/user/:userID/:ID", baseInstance.DeleteExitPermission)
	}
}
