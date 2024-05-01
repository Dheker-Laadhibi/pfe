package questions

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// CreateQuestion 	Handles the creation of a new question.
// @Summary        	Create question
// @Description    	Create a new question.
// @Tags			Questions
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				  true		"Company ID"
// @Param			request			body			questions.QuestionsIn true		"Question query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	  "Invalid request"
// @Failure			401				{object}		utils.ApiResponses	  "Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	  "Forbidden"
// @Failure			500				{object}		utils.ApiResponses	  "Internal Server Error"
// @Router			/questions/{companyID}	[post]
func (db Database) CreateQuestion(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a QuestionIn struct
	question := new(QuestionsIn)
	if err := ctx.ShouldBindJSON(question); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Convert options slice to pq.StringArray
	optionsArray := pq.StringArray(question.Options)

	// Create a new question in the database
	dbQuestion := &domains.Questions{
		ID:                   uuid.New(),
		Question:             question.Question,
		CorrectAnswer:        question.CorrectAnswer,
		Options:              optionsArray,
		AssociatedTechnology: question.AssociatedTechnology,
		CompanyID:            session.CompanyID,
	}

	if err := domains.Create(db.DB, dbQuestion); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadQuestions 	Handles the retrieval of all questions.
// @Summary        	Get questions
// @Description    	Get all questions.
// @Tags			Questions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page				query		int				false		"Page"
// @Param			limit				query		int				false		"Limit"
// @Param			companyID			path		string			true		"Company ID"
// @Success			200					{object}	questions.QuestionsPagination
// @Failure			400					{object}	utils.ApiResponses			"Invalid request"
// @Failure			401					{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}	utils.ApiResponses			"Forbidden"
// @Failure			500					{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/questions/{companyID}	[get]
func (db Database) ReadQuestions(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the page from the request parameter
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(constants.DEFAULT_PAGE_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the limit from the request parameter
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(constants.DEFAULT_LIMIT_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the question's value is among the allowed choices
	validChoices := utils.ResponseLimitPagination()
	isValidChoice := false
	for _, choice := range validChoices {
		if uint(limit) == choice {
			isValidChoice = true
			break
		}
	}

	// If the value is invalid, set it to default DEFAULT_LIMIT_PAGINATION
	if !isValidChoice {
		limit = constants.DEFAULT_LIMIT_PAGINATION
	}

	// Generate offset
	offset := (page - 1) * limit

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all question data from the database
	questions, err := ReadAllPagination(db.DB, []domains.Questions{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all questions data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Questions{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a question structure as a response
	response := QuestionsPagination{}
	dataTableQuestion := []QuestionsTable{}
	for _, question := range questions {

		dataTableQuestion = append(dataTableQuestion, QuestionsTable{
			ID:                   question.ID,
			Question:             question.Question,
			CorrectAnswer:        question.CorrectAnswer,
			Options:              question.Options,
			AssociatedTechnology: question.AssociatedTechnology,
			CreatedAt:            question.CreatedAt,
		})
	}
	response.Items = dataTableQuestion
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadQuestionsList Handles the retrieval the list of all questions.
// @Summary        	Get list of  questions
// @Description    	Get list of all questions.
// @Tags			Questions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Success			200						{array}			questions.QuestionsList
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/questions/{companyID}/list	[get]
func (db Database) ReadQuestionsList(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all question data from the database
	questions, err := ReadAllList(db.DB, []domains.Questions{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all question data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a question structure as a response
	questionsList := []QuestionsList{}
	for _, question := range questions {

		questionsList = append(questionsList, QuestionsList{
			ID:                   question.ID,
			Question:             question.Question,
			CorrectAnswer:        question.CorrectAnswer,
			Options:              question.Options,
			AssociatedTechnology: question.AssociatedTechnology,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, questionsList)
}

// ReadQuestionsCount	Handles the retrieval the number of all questions.
// @Summary        	Get number of  questions
// @Description    	Get number of all questions.
// @Tags			Questions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true		"Company ID"
// @Success			200						{object}		questions.QuestionsCount
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/questions/{companyID}/count	[get]
func (db Database) ReadQuestionsCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all question data from the database
	questions, err := domains.ReadTotalCount(db.DB, &[]domains.Questions{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all question data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a question structure as a response
	questionsCount := QuestionsCount{
		Count: questions,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, questionsCount)
}

// ReadQuestion 	Handles the retrieval of one question.
// @Summary        	Get question
// @Description    	Get one question.
// @Tags			Questions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Param			ID						path			string			true	"Question ID"
// @Success			200						{object}		questions.QuestionsDetails
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/questions/{companyID}/{ID}	[get]
func (db Database) ReadQuestion(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the question ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve question data from the database
	question, err := ReadByID(db.DB, domains.Questions{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving question data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Retriece name from the database
	companyName, err := domains.ReadCompanyNameByID(db.DB, session.CompanyID)
	if err != nil {
		logrus.Error("Error retrieving company name data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a question structure as a response
	details := QuestionsDetails{
		ID:                   question.ID,
		Question:             question.Question,
		CorrectAnswer:        question.CorrectAnswer,
		Options:              question.Options,
		AssociatedTechnology: question.AssociatedTechnology,
		CompanyName:          companyName,
		CreatedAt:            question.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateQuestion 	Handles the update of a question.
// @Summary        	Update question
// @Description    	Update one question.
// @Tags			Questions
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true		"Company ID"
// @Param			ID					path			string				true		"Question ID"
// @Param			request				body			questions.QuestionsIn		true		"Question query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses				"Invalid request"
// @Failure			401					{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses				"Forbidden"
// @Failure			500					{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/questions/{companyID}/{ID}	[put]
func (db Database) UpdateQuestion(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the question ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a QuestionsIn struct
	question := new(QuestionsIn)
	if err := ctx.ShouldBindJSON(question); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the question with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Questions{}, objectID); err != nil {
		logrus.Error("Error checking if the question with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the question data in the database
	dbQuestion := &domains.Questions{
		Question:             question.Question,
		CorrectAnswer:        question.CorrectAnswer,
		Options:              question.Options,
		AssociatedTechnology: question.AssociatedTechnology,
		CompanyID:            session.CompanyID,
	}
	if err = domains.Update(db.DB, dbQuestion, objectID); err != nil {
		logrus.Error("Error updating question data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteQuestion	Handles the deletion of a question.
// @Summary        	Delete question
// @Description    	Delete one question.
// @Tags			Questions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true		"Company ID"
// @Param			ID					path			string			true		"Question ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/questions/{companyID}/{ID}	[delete]
func (db Database) DeleteQuestion(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the question ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the question with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Questions{}, objectID); err != nil {
		logrus.Error("Error checking if the question with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the question data from the database
	if err := domains.Delete(db.DB, &domains.Questions{}, objectID); err != nil {
		logrus.Error("Error deleting question data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
