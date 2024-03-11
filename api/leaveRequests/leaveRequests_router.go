package leaveRequests

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LeaveRouterInit initializes the routes related to leave.
func LeaveRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewLeaveRequestsRepository(db)

	// Private
	leave := router.Group("/LeaveRequests")
	{

		//Create a new Leave demande
		leave.POST("", baseInstance.AddLeave)

		// GET endpoint to retrieve all leave for a specific user
		leave.GET("/:userID", baseInstance.ReadLeave)

		// GET endpoint to retrieve the count of leave for a specific user
		leave.GET("/:userID/count", baseInstance.ReadLeaveCount)

		// GET endpoint to retrieve details of a specific leave for a specific user
		leave.GET("/:userID/:ID", baseInstance.ReadOneLeave)

		// PUT endpoint to update the details of a specific leave for a specific user
		leave.PUT("/:userID/:ID", baseInstance.UpdateLeave)

		// DELETE endpoint to delete a specific leave for a specific user
		leave.DELETE("/:userID/:ID", baseInstance.DeleteLeave)
	}

}
