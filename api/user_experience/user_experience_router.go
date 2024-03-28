package user_experience

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRouterInit initializes the routes related to projects.
func ExperienceRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewExperienceRepository(db)

	// Private
	projects := router.Group("/experience/:companyID")
	{

		// POST endpoint to create a new project
		projects.POST("/create/:userID", baseInstance.CreateUserExperience)

		// GET endpoint to retrieve all projects for a specific company
		projects.GET("", baseInstance.ReadUserExperience)

		// GET endpoint to retrieve a list of projects for a specific company
		projects.GET("/list", baseInstance.ReadExperiencesList)

		// GET endpoint to retrieve the count of projects for a specific company
		projects.GET("/count", baseInstance.ReadUserExperiencesCount)

		// GET endpoint to retrieve details of a specific project
		projects.GET("/:ID/Get", baseInstance.ReadOneUserExperiences)

		// PUT endpoint to update details of a specific project
		projects.PUT("/:ID/update", baseInstance.UpdateUserExperience)

		// DELETE endpoint to delete a specific project
		projects.DELETE("/:ID", baseInstance.DeleteExperience)

	}
}
