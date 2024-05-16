package permissions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InternRouterInit initializes the routes related to interns.
func PermissionRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewPermissionRepository(db)


	/**

	IMPORTANT:
	The user ID represents the unique identifier of an employee who holds the role of former intern groups.
	Please ensure that appropriate permissions and access controls are in place for this user.
	*/
	// Private
	permissions := router.Group("/permissions/:companyID")
	{

		// POST endpoint to create a new intern
		permissions.POST("/:featureID/:roleID", baseInstance.CreatePermission)

		// GET endpoint to retrieve all interns for a specific company and a specific user
		permissions.GET("", baseInstance.ReadPermissions)

		// GET endpoint to retrieve the count of interns for a specific company
		permissions.GET("/count", baseInstance.ReadPermissionCount)

		// GET endpoint to retrieve details of a specific intern
		permissions.GET("/permission/:permissionID", baseInstance.ReadPermission)

		// PUT endpoint to update details of a specific intern
		permissions.PUT("/update/:permissionID", baseInstance.UpdatePermission)

		// DELETE endpoint to delete a specific intern
		permissions.DELETE("/delete/:permissionID", baseInstance.DeletePermission)


	}

}
