package projects

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRouterInit initializes the routes related to projects.
func ProjectRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewProjectRepository(db)

	// Private
	projects := router.Group("/projects/:companyID")
	{

		// POST endpoint to create a new project
		projects.POST("", baseInstance.CreateProject)

		// GET endpoint to retrieve all projects for a specific company
		projects.GET("", baseInstance.ReadProjects)

		// GET endpoint to retrieve a list of projects for a specific company
		projects.GET("/list", baseInstance.ReadProjectsList)

		// GET endpoint to retrieve the count of projects for a specific company
		projects.GET("/count", baseInstance.ReadProjectsCount)

		// GET endpoint to retrieve details of a specific project
		projects.GET("/:ID", baseInstance.ReadProject)

		// PUT endpoint to update details of a specific project
		projects.PUT("/:ID", baseInstance.UpdateProject)

		// DELETE endpoint to delete a specific project
		projects.DELETE("/:ID", baseInstance.DeleteProject)
	}
}
