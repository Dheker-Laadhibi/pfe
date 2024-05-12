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
	interns := router.Group("/interns/:companyID")
	{

		// POST endpoint to create a new intern
		interns.POST("/add/:supervisorID", baseInstance.CreatePermission)

		// GET endpoint to retrieve all interns for a specific company and a specific user 
		interns.GET("", baseInstance.ReadPermissions)
		// GET endpoint to retrieve the count of interns for a specific company
		interns.GET("/count", baseInstance.ReadPermissionCount)

		// GET endpoint to retrieve details of a specific intern
		interns.GET("/intern/:internID", baseInstance.ReadPermission)

		// PUT endpoint to update details of a specific intern
		interns.PUT("/update/:internID", baseInstance.UpdatePermission)

		// DELETE endpoint to delete a specific intern
		interns.DELETE("/delete/:internID", baseInstance.DeletePermission)
	}
}
