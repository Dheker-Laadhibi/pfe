package interns

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InternRouterInit initializes the routes related to interns.
func InternRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewInternRepository(db)

	/**

IMPORTANT:
The user ID represents the unique identifier of an employee who holds the role of former intern groups.
Please ensure that appropriate permissions and access controls are in place for this user.
*/
	// Private
	interns := router.Group("/interns")
	{

		// POST endpoint to create a new intern
		interns.POST("/:companyID", baseInstance.CreateIntern)

		// GET endpoint to retrieve all interns for a specific company and a specific user 
		interns.GET("all/:companyID/:SupervisorID", baseInstance.ReadInterns)
		// GET endpoint to retrieve the count of interns for a specific company
		interns.GET("/count/:companyID", baseInstance.ReadInternsCount)

		// GET endpoint to retrieve details of a specific intern
		interns.GET("/:companyID/:ID", baseInstance.ReadIntern)

		// PUT endpoint to update details of a specific intern
		interns.PUT("update/:companyID/:ID", baseInstance.UpdateIntern)

		// DELETE endpoint to delete a specific intern
		interns.DELETE("Delete/:companyID/:SupervisorID", baseInstance.DeleteIntern)
	}
}
