package loanRequests

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoanRequestRouterInit initializes the routes related to loanRequests.
func LoanRequestRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewLoanRequestRepository(db)

	// Private
	loanRequests := router.Group("/loanRequests")
	{

		//Create a new loanRequests demande
		loanRequests.POST("", baseInstance.AddLoanRequests)

		// GET endpoint to retrieve all loanRequests for a specific company
		loanRequests.GET("/:userID", baseInstance.ReadAllLoanRequests)

		// GET endpoint to retrieve the count of loanRequests for a specific user
		loanRequests.GET("/user/:userID/count", baseInstance.ReadLoanRequestsCount)

		// GET endpoint to retrieve details of a specific loanRequests for a specific user
		loanRequests.GET("/user/:userID/:ID", baseInstance.ReadOneLoanRequests)

		// PUT endpoint to update the details of a specific loanRequests for a specific user
		loanRequests.PUT("/user/:userID/:ID", baseInstance.UpdateLoanRequests)

		// DELETE endpoint to delete a specific loanRequests for a specific user
		loanRequests.DELETE("/user/:userID/:ID", baseInstance.DeleteLoanRequests)
	}
}
