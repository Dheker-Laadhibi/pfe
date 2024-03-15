package mission_orders

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PresenceRouterInit initializes the routes related to presences.
func MissionOrdersRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewNewTrainingRequestRepository(db)

	// Private
	MissionOrders := router.Group("/missions")
	{

		// create endpoint to create a specific MissionOrders for a specific user
		MissionOrders.POST("/:companyID", baseInstance.CreateTrainingRequestByUser)

		// GET endpoint to retrieve all MissionOrders for a specific user
		MissionOrders.GET("/All/:userID", baseInstance.ReadMissionsOrders)

		// GET endpoint to retrieve the count of MissionOrders for a specific user
		MissionOrders.GET("/count/:userID", baseInstance.ReadMissionOrdersCount)

		// GET endpoint to retrieve details of a specific MissionOrders for a specific user
		MissionOrders.GET("/get/:ID/:userID", baseInstance.ReadMissionOrders)

		// PUT endpoint to update the details of a specific MissionOrders for a specific user
		MissionOrders.PUT("/update/:ID/:userID", baseInstance.UpdateMissionOrders)

		// DELETE endpoint to delete a specific MissionsOrders for a specific user
		MissionOrders.DELETE("/delete/:ID/:userID", baseInstance.DeleteMissionOrders)




	}
}
