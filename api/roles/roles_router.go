package roles

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleRouterInit initializes the routes related to roles.
func RoleRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewRoleRepository(db)

	// Private
	roles := router.Group("/roles/:companyID")
	{

		// POST endpoint to create a new role
		roles.POST("", baseInstance.CreateRole)

		// GET endpoint to retrieve all roles for a specific company
		roles.GET("", baseInstance.ReadRoles)

		// GET endpoint to retrieve a list of roles for a specific company
		roles.GET("/list", baseInstance.ReadRolesList)

		// GET endpoint to retrieve the count of roles for a specific company
		roles.GET("/count", baseInstance.ReadRolesCount)

		// GET endpoint to retrieve details of a specific role
		roles.GET("/:ID", baseInstance.ReadRole)

		// PUT endpoint to update a specific role
		roles.PUT("/:ID", baseInstance.UpdateRole)

		// DELETE endpoint to delete a specific role
		roles.DELETE("/:ID", baseInstance.DeleteRole)
	}
}
