package presences

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PresenceRouterInit initializes the routes related to presences.
func PresenceRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewPresenceRepository(db)

	// Private
	presences := router.Group("/presences")
	{

		// create endpoint to crete a specific presences for a specific user
		presences.POST("/:companyID", baseInstance.CreatePresence)

		// GET endpoint to retrieve all presences for a specific user
		presences.GET("/All/:userID", baseInstance.ReadPresences)

		// GET endpoint to retrieve the count of presences for a specific user
		presences.GET("/count/:userID", baseInstance.ReadPresencesCount)

		// GET endpoint to retrieve details of a specific presences for a specific user
		presences.GET("/get/:ID/:userID", baseInstance.ReadPresence)

		// PUT endpoint to update the details of a specific presences for a specific user
		presences.PUT("/update/:ID/:userID", baseInstance.UpdatePresence)

		// DELETE endpoint to delete a specific presences for a specific user
		presences.DELETE("/delete/:ID/:userID", baseInstance.DeletePresence)
	}
}
