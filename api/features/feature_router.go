package features

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleRouterInit initializes the routes related to roles.
func FeatureRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewFeatureRepository(db)

	// Private
	features := router.Group("/features/:companyID")
	{

		// POST endpoint to create a new feature
		features.POST("", baseInstance.CreateFeature)

		// GET endpoint to retrieve all feature for a specific company
		features.GET("", baseInstance.ReadFeatures )

		// GET endpoint to retrieve a list of feature for a specific company
		features.GET("/list", baseInstance.ReadFeaturesList)

		// GET endpoint to retrieve the count of feature for a specific company
		features.GET("/count", baseInstance.ReadFeaturesCount)

		// GET endpoint to retrieve details of a specific feature
		features.GET("/:ID", baseInstance.ReadFeature)

		// PUT endpoint to update a specific feature
		features.PUT("/:ID", baseInstance.UpdateFeature)

		// DELETE endpoint to delete a specific feature
		features.DELETE("/:ID", baseInstance.DeleteFeature)

	}
}
