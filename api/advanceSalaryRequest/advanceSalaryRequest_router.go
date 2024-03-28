package advanceSalaryRequest

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdvanceSalaryRequestRouterInit initializes the routes related to ladvanceSalaryRequest.
func AdvanceSalaryRequestsRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewAdvanceSalaryRequestsRepository(db)

	// Private
	AdvanceSalaryRequest := router.Group("/advance_salaryRequests")
	{
		//Create a new AdvanceSalaryRequest demande
		AdvanceSalaryRequest.POST("/:companyID", baseInstance.AddAdvanceSalaryRequest)

		// GET endpoint to retrieve all AdvanceSalaryRequest for a specific user   advance_salaryRequests/company/{userID}
		AdvanceSalaryRequest.GET("/user/:userID", baseInstance.ReadAdvanceSalaryRequest)

		// GET endpoint to retrieve all AdvanceSalaryRequest for a specific company
		AdvanceSalaryRequest.GET("/company/:companyID", baseInstance.ReadAdvanceSalaryRequestsByCompany)

		// GET endpoint to retrieve the count of AdvanceSalaryRequest for a specific user
		AdvanceSalaryRequest.GET("/:userID/count", baseInstance.ReadAdvanceSalaryRequestCount)

		// GET endpoint to retrieve the count of AdvanceSalaryRequest for a specific company
		AdvanceSalaryRequest.GET("/count/:companyID", baseInstance.ReadAdvanceSalaryRequestCountbycompany)

		// GET endpoint to retrieve details of a specific AdvanceSalaryRequest for a specific user
		AdvanceSalaryRequest.GET("/:userID/:ID", baseInstance.ReadOneAdvanceSalaryRequest)

		// PUT endpoint to update the details of a specific AdvanceSalaryRequest for a specific user
		AdvanceSalaryRequest.PUT("/:userID/:ID", baseInstance.UpdateAdvanceSalaryRequest)

		// DELETE endpoint to delete a specific AdvanceSalaryRequest for a specific user
		AdvanceSalaryRequest.DELETE("/:userID/:ID", baseInstance.DeleteAdvanceSalaryRequest)
	}

}
