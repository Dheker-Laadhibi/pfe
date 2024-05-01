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
	TrainingRequest := router.Group("/training_request")
	{

		// create endpoint to create a specific trainings for a specific company
		TrainingRequest.POST("/:companyID/:userID", baseInstance.CreateTrainingRequestByUser)

		// GET endpoint to retrieve all trainings for a specific company
		TrainingRequest.GET("/All/:companyID", baseInstance.ReadTrainingsRequest)

		// GET endpoint to retrieve the count of trainings  for a specific company
		TrainingRequest.GET("/count/:companyID", baseInstance.ReadTrainingsCount)

		// GET endpoint to retrieve details of a specific trainings for a specific company
		TrainingRequest.GET("/get/:companyID/:ID", baseInstance.ReadTrainingsRequests)

		// PUT endpoint to update the details of a specific training for a specific company
		TrainingRequest.PUT("/update/:companyID/:ID", baseInstance.UpdateTraining)

		// DELETE endpoint to delete a specific training for a specific company
		TrainingRequest.DELETE("/delete/:companyID/:ID", baseInstance.DeleteTrainingsRequest)

	}
}
