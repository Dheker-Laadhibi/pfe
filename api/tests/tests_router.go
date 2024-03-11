package tests

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRouterInit initializes the routes related to tests.
func TestRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewTestRepository(db)

	// Private
	tests := router.Group("/tests/:companyID")
	{

		// POST endpoint to create a new test
		tests.POST("", baseInstance.CreateTest)

		// POST endpoint to Generate Questions
		tests.POST("/:ID", baseInstance.GenerateQuestions)

		// POST endpoint to Assign Tests
		tests.POST("/Assign/:candidatID", baseInstance.AssignTest)

		// GET endpoint to retrieve all tests for a specific company
		tests.GET("", baseInstance.ReadTests)

		// GET endpoint to retrieve all candidats for a specific test
		tests.GET("/Candidats/:testID", baseInstance.ReadTestsCandidats)

		// GET endpoint to retrieve all questionss for a specific test
		tests.GET("/Questions/:testID", baseInstance.ReadTestsQuestions)

		// GET endpoint to retrieve a list of tests for a specific company
		tests.GET("/list", baseInstance.ReadTestsList)

		// GET endpoint to retrieve the count of tests for a specific company
		tests.GET("/count", baseInstance.ReadTestsCount)

		// GET endpoint to retrieve details of a specific test
		tests.GET("/:ID", baseInstance.ReadTest)

		// PUT endpoint to update details of a specific test
		tests.PUT("/:ID", baseInstance.UpdateTest)

		// DELETE endpoint to delete a specific test
		tests.DELETE("/:ID", baseInstance.DeleteTest)
	}
}
