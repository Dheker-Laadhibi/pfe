package tests

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateTest 		Handles the creation of a new test.
// @Summary        	Create test
// @Description    	Create a new test.
// @Tags			Tests
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				true		"Company ID"
// @Param			candidatID		path			string				true		"Candidat ID"
// @Param			nbrQuestions		query			string				false		"NbrQuestions"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/tests/{companyID}/create/{candidatID}	[post]
func (db Database) CreateTest(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the page from the request parameter
	nbrQuestions, err := strconv.Atoi(ctx.DefaultQuery("nbrQuestions", strconv.Itoa(constants.DEFAULT_PAGE_PAGINATION)))
	if err != nil {
		// Handle invalid integer format
		logrus.Error(ctx, http.StatusBadRequest, "Invalid number of questions", err)
		return
	}

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		// Handle invalid UUID format
		logrus.Error(ctx, http.StatusBadRequest, "Invalid company ID", err)
		return
	}

	// Parse and validate the company ID from the request parameter
	candidatID, err := uuid.Parse(ctx.Param("candidatID"))
	if err != nil {
		// Handle invalid UUID format
		logrus.Error(ctx, http.StatusBadRequest, "Invalid candidat ID", err)
		return
	}

	

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		// Handle employee verification error
		logrus.Error(ctx, http.StatusBadRequest, "Error verifying employee belonging", err)
		return
	}

	//Retrive the projectID based on the candidatID
	projectID, err := ReadProjectIDByCandidatID(db.DB, candidatID)
	if err != nil {
		// Handle error retrieving projectID
		logrus.Error(ctx, http.StatusBadRequest, "Error retrieving projectID", err)
		return
	}
	
	//Retrive the choosen project details
	project, err := ReadProjectDetails(db.DB, domains.Project{}, projectID)
	if err != nil {
		// Handle error retrieving project details
		logrus.Error(ctx, http.StatusBadRequest, "Error retrieving project details", err)
		return
	}

	//generate a new uuid "test_id"
	var ID = uuid.New()
	var TestName = project.Projectname + "Test"

	// Create a new test in the database
	dbTest := &domains.Tests{
		ID:           ID,
		Title:        TestName,
		Specialty:    project.Specialty,
		Technologies: project.Technologies,
		CompanyID:    session.CompanyID,
	}
	if err := domains.Create(db.DB, dbTest); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Create a test structure as a response
	dbTestCandidat := &domains.TestCandidats{
		TestID:     ID,
		CandidatID: candidatID,
		CompanyID:  companyID,
	}
	if err := domains.Create(db.DB, dbTestCandidat); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve the list of technologies associated with the test
	technologies, err := ReadTechList(db.DB, ID)
	if err != nil {
		// Handle error retrieving technologies
		logrus.Error(ctx, http.StatusBadRequest, "Error retrieving technologies for the test", err)
		return
	}

	// Retrieve the list of questions associated with the test technologies
	questions, err := GetAllQuestions(db.DB, technologies)
	if err != nil {
		// Handle error retrieving questions
		logrus.Error(ctx, http.StatusBadRequest, "Error retrieving questions for the test", err)
		return
	}

	// Generate a random selection of questions
	selectedQuestions, err := GetRandomQuestions(questions, nbrQuestions)
	if err != nil {
		// Handle error generating random questions
		logrus.Error(ctx, http.StatusBadRequest, "Error generating random questions", err)
		return
	}

	for _, selectedQuestion := range selectedQuestions {
		Qest, _ := uuid.Parse(selectedQuestion)
		question, err := ReadQuestionDetails(db.DB, domains.Questions{}, Qest)
		if err != nil {
			// Handle error retrieving question details
			logrus.Error(ctx, http.StatusBadRequest, "Error retrieving question details", err)
			return
		}
		// Create a new entry for the test with the selected questions
		dbTestQuestions := &domains.TestQuestions{
			TestID:               ID,
			QuestionID:           question.ID,
			Question:             question.Question,
			CorrectAnswer:        question.CorrectAnswer,
			Options:              question.Options,
			AssociatedTechnology: question.AssociatedTechnology,
			CandidatID:           candidatID,
			CompanyID:            session.CompanyID,
		}
		if err := domains.Create(db.DB, dbTestQuestions); err != nil {
			// Handle error saving data to the database
			logrus.Error(ctx, http.StatusBadRequest, "Error saving data to the database", err)
			return
		}

	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadTests 		Handles the retrieval of all tests.
// @Summary        	Get tests
// @Description    	Get all tests.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		path		string		true		"Company ID"
// @Success			200				{object}	tests.TestsPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/tests/{companyID}	[get]
func (db Database) ReadTests(ctx *gin.Context) {

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

	// Check if the test's value is among the allowed choices
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

	// Retrieve all test data from the database
	tests, err := ReadAllPagination(db.DB, []domains.Tests{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all test data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Tests{}, "company_id", companyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a test structure as a response
	response := TestsPagination{}
	dataTableTest := []TestsTable{}
	for _, test := range tests {

		dataTableTest = append(dataTableTest, TestsTable{
			ID:           test.ID,
			Title:        test.Title,
			Specialty:    test.Specialty,
			Technologies: test.Technologies,
			CreatedAt:    test.CreatedAt,
		})
	}
	response.Items = dataTableTest
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

//***************************************************

// ReadQuestionsbyTest		Handles the retrieval of all Questions by testID.
// @Summary        	Get Questions
// @Description    	Get all Questions by testID.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		path		string		true		"Test ID"
// @Param			testID		path		string		true		"Test ID"
// @Success			200				{object}	tests.QuestionsPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/tests/{companyID}/TQuestions/{testID}	[get]
func (db Database) ReadQuestionsbyTest(ctx *gin.Context) {

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
	testID, err := uuid.Parse(ctx.Param("testID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the test's value is among the allowed choices
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

	// Retrieve all test data from the database
	questions, err := ReadAllPaginationQT(db.DB, []domains.TestQuestions{}, testID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all test data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCountQT(db.DB, &domains.TestQuestions{}, "tests_id", testID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a test structure as a response
	response := QuestionsPagination{}
	dataTableQuestion := []QuestionsTable{}
	for _, Question := range questions {

		dataTableQuestion = append(dataTableQuestion, QuestionsTable{
			QuestionID: Question.QuestionID,
			Question:   Question.Question,
			Options:    Question.Options,
		})
	}
	response.Items = dataTableQuestion
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

//***************************************************

// ReadScores 		Handles the retrieval of all scores.
// @Summary        	Get candidats scores
// @Description    	Get all scores.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		path		string		true		"Company ID"
// @Success			200				{object}	tests.ScoresPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/tests/{companyID}/scores	[get]
func (db Database) ReadScores(ctx *gin.Context) {

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

	// Check if the test's value is among the allowed choices
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

	// Retrieve all test data from the database
	scores, err := ReadAllPaginationS(db.DB, []domains.TestCandidats{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all test data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCountAS(db.DB, &domains.TestCandidats{}, "company_id", companyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a test structure as a response
	response := ScoresPagination{}
	dataTableScore := []ScorsTable{}
	for _, score := range scores {
		// Retrieve test data from the database
		test, err := ReadByID(db.DB, domains.Tests{}, score.TestID)
		if err != nil {
			logrus.Error("Error occurred while finding all test data. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

		// Retrieve candidat data from the database
		candidat, err := ReadCandidatByID(db.DB, domains.Condidats{}, score.CandidatID)
		if err != nil {
			logrus.Error("Error occurred while finding all test data. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

		dataTableScore = append(dataTableScore, ScorsTable{
			IdTest:     score.TestID,
			IdCandidat: score.CandidatID,
			Title:      test.Title,
			Firstname:  candidat.Firstname,
			Lastname:   candidat.Lastname,
			Score:      score.Score,
		})
	}
	response.Items = dataTableScore
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadTestsList 	Handles the retrieval the list of all tests.
// @Summary        	Get list of  tests
// @Description    	Get list of all tests.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Success			200					{array}			tests.TestsList
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/tests/{companyID}/list	[get]
func (db Database) ReadTestsList(ctx *gin.Context) {

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

	// Retrieve all test data from the database
	tests, err := ReadAllList(db.DB, []domains.Tests{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all test data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a test structure as a response
	usersList := []TestsList{}
	for _, test := range tests {
		usersList = append(usersList, TestsList{
			ID:           test.ID,
			Specialty:    test.Specialty,
			Technologies: test.Technologies,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, usersList)
}

// ReadTestsCount 	Handles the retrieval the number of all tests.
// @Summary        	Get number of  tests
// @Description    	Get number of all tests.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true	"Company ID"
// @Success			200						{object}		tests.TestsCount
// @Failure			400						{object}		utils.ApiResponses	"Invalid request"
// @Failure			401						{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses	"Forbidden"
// @Failure			500						{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/tests/{companyID}/count	[get]
func (db Database) ReadTestsCount(ctx *gin.Context) {

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

	// Retrieve all test data from the database
	tests, err := domains.ReadTotalCount(db.DB, &[]domains.Tests{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all test data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a test structure as a response
	TestsCount := TestsCount{
		Count: tests,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, TestsCount)
}

// UpdateCandidatAnswer 		Handles the update the candidat response.
// @Summary        	Update the candidat response
// @Description    	Update the candidat response for a specific question.
// @Tags			Tests
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true	"Company ID"
// @Param			candidatID			path			string				true	"Candidat ID"
// @Param			testID			path			string				true	"Test ID"
// @Param			questionID			path			string				true	"Question ID"
// @Param			request				body			tests.CandidatAnswerIn		true	"CandidatAnswerIn query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/tests/{companyID}/{candidatID}/{testID}/{questionID}	[put]
func (db Database) UpdateCandidatAnswer(ctx *gin.Context) {

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

	// Parse and validate the question ID from the request parameter
	questionID, err := uuid.Parse(ctx.Param("questionID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the candidat ID from the request parameter
	candidatID, err := uuid.Parse(ctx.Param("candidatID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the test ID from the request parameter
	testID, err := uuid.Parse(ctx.Param("testID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	//check if the user is a candidat by the session.role_id

	// Parse the incoming JSON request into a UserIn struct
	response := new(CandidatAnswer)
	if err := ctx.ShouldBindJSON(response); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the test data in the database
	dbTest := &domains.TestQuestions{
		CandidatAnswer: response.CandidatAnswer,
	}
	if err = domains.UpdateAnswer(db.DB, dbTest, questionID, candidatID, testID); err != nil {
		logrus.Error("Error updating test data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteTest	 	Handles the deletion of a test.
// @Summary        	Delete test
// @Description    	Delete one test.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"Test ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/tests/{companyID}/{ID}	[delete]
func (db Database) DeleteTest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the test ID from the request parameter
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

	// Check if the test with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tests{}, objectID); err != nil {
		logrus.Error("Error checking if the test with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the test data from the database
	if err := domains.Delete(db.DB, &domains.Tests{}, objectID); err != nil {
		logrus.Error("Error deleting test data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// ReadTestAnswers 		Handles the retrieval of all the answers of a test by candidatID and Insert the score of the candidat into the database.
// @Summary        	Get candidat Responses
// @Description    	Get  all the answers of a test by candidatID.
// @Tags			Tests
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID			path			string			true	"Company ID"
// @Param			candidatID					path			string			true	"Candidat ID"
// @Param			testID					path			string			true	"Test ID"
// @Success			200					{object}		tests.TestsDetails
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/tests/{companyID}/{candidatID}/{testID}	[get]
func (db Database) ReadTestAnswers(ctx *gin.Context) {

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

	// Parse and validate the candidat ID from the request parameter
	candidatID, err := uuid.Parse(ctx.Param("candidatID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the test ID from the request parameter
	testID, err := uuid.Parse(ctx.Param("testID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the test's value is among the allowed choices
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

	// Retrieve all test data from the database
	answers, err := ReadAllPaginationAnswers(db.DB, []domains.TestQuestions{}, testID, candidatID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all test data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCountAnswers(db.DB, &domains.TestQuestions{}, "tests_id", "candidat_id", testID, candidatID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a test structure as a response
	var Score = 0
	response := AnswersPagination{}
	dataTableQuestion := []CandidatAnswer{}
	for _, Answers := range answers {
		if Answers.CorrectAnswer == Answers.CandidatAnswer {
			Score = Score + 1
		}
		dataTableQuestion = append(dataTableQuestion, CandidatAnswer{
			QuestionID:     Answers.QuestionID,
			Question:       Answers.Question,
			CorrectAnswer:  Answers.CorrectAnswer,
			CandidatAnswer: Answers.CandidatAnswer,
		})
	}
	response.Items = dataTableQuestion
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	if err = domains.UpdateScore(db.DB, &domains.TestCandidats{}, candidatID, testID, Score); err != nil {
		logrus.Error("Error updating test data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}
