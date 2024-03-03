package condidats

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleRouterInit initializes the routes related to condidats.
func CondidatRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewCondidatRepository(db)

	// Private
	condidats := router.Group("/condidats/:companyID")
	{

		// POST endpoint to create a new condidat
		condidats.POST("", baseInstance.CreateCondidat)

		// GET endpoint to retrieve all condidats for a specific company
		condidats.GET("", baseInstance.ReadCondidats)

		// GET endpoint to retrieve a list of condidats for a specific company
		condidats.GET("/list", baseInstance.ReadCondidatsList)

		// GET endpoint to retrieve the count of condidats for a specific company
		condidats.GET("/count", baseInstance.ReadCondidatsCount)

		// GET endpoint to retrieve details of a specific condidat
		condidats.GET("/:ID", baseInstance.ReadCondidat)

		// PUT endpoint to update a specific condidat
		condidats.PUT("/:ID", baseInstance.UpdateCondidat)

		// DELETE endpoint to delete a specific condidat
		condidats.DELETE("/:ID", baseInstance.DeleteCondidat)
	}
}
