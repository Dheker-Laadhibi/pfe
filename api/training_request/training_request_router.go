package training_request

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TrainingRequestRouterInit initializes the routes related to TrainingRequest.
func TrainingRequestRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewNewTrainingRequestRepository(db)

	// Private
	TrainingRequest := router.Group("/missions")
	{

		// create endpoint to create a specific MissionOrders for a specific user
		TrainingRequest.POST("/:companyID", baseInstance.CreateTrainingRequestByUser)

		// GET endpoint to retrieve all MissionOrders for a specific user
		TrainingRequest.GET("/All/:userID", baseInstance.ReadMissionsOrders)

		// GET endpoint to retrieve the count of MissionOrders for a specific user
		TrainingRequest.GET("/count/:userID", baseInstance.ReadMissionOrdersCount)

		// GET endpoint to retrieve details of a specific MissionOrders for a specific user
		TrainingRequest.GET("/get/:ID/:userID", baseInstance.ReadMissionOrders)

		// PUT endpoint to update the details of a specific MissionOrders for a specific user
		TrainingRequest.PUT("/update/:ID/:userID", baseInstance.UpdateMissionOrders)

		// DELETE endpoint to delete a specific MissionsOrders for a specific user
		TrainingRequest.DELETE("/delete/:ID/:userID", baseInstance.DeleteMissionOrders)

	}
}
