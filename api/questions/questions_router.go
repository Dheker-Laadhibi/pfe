package questions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleRouterInit initializes the routes related to questions.
func QuestionRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewQuestionRepository(db)

	// Private
	questions := router.Group("/questions/:companyID")
	{

		// POST endpoint to create a new question
		questions.POST("", baseInstance.CreateQuestion)

		// GET endpoint to retrieve all questions for a specific company
		questions.GET("", baseInstance.ReadQuestions)

		// GET endpoint to retrieve a list of questions for a specific company
		questions.GET("/list", baseInstance.ReadQuestionsList)

		// GET endpoint to retrieve the count of questions for a specific company
		questions.GET("/count", baseInstance.ReadQuestionsCount)

		// GET endpoint to retrieve details of a specific question
		questions.GET("/:ID", baseInstance.ReadQuestion)

		// PUT endpoint to update a specific question
		questions.PUT("/:ID", baseInstance.UpdateQuestion)

		// DELETE endpoint to delete a specific question
		questions.DELETE("/:ID", baseInstance.DeleteQuestion)
	}
}
