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
		tests.POST("/create/:candidatID", baseInstance.CreateTest)

		// GET endpoint to retrieve all tests for a specific company
		tests.GET("", baseInstance.ReadTests)

		// GET endpoint to retrieve all questions for a specific Test
		tests.GET("/TQuestions/:testID", baseInstance.ReadQuestionsbyTest)

		// GET endpoint to retrieve all Answers details for a specific test and a specific candidat
		tests.GET("/:candidatID/:testID", baseInstance.ReadTestAnswers)

		// GET endpoint to retrieve all candidats for a specific test
		tests.GET("/scores", baseInstance.ReadScores)

		// GET endpoint to retrieve a list of tests for a specific company
		tests.GET("/list", baseInstance.ReadTestsList)

		// GET endpoint to retrieve the count of tests for a specific company
		tests.GET("/count", baseInstance.ReadTestsCount)

		// PUT endpoint to Update the candidat answer
		tests.PUT("/:candidatID/:testID/:questionID", baseInstance.UpdateCandidatAnswer)

		// DELETE endpoint to delete a specific test
		tests.DELETE("/:ID", baseInstance.DeleteTest)
	}
}
