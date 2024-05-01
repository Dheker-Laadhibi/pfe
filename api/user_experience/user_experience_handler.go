package user_experience

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

// CreateUserExperience	Handles the creation of a new UserExperience	.
// @Summary        		Create UserExperience
// @Description    		Create a new UserExperience.
// @Tags				UserExperience
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID		path				string				true			"companyID"
// @Param				userID			path			    string				true		    "userID"
// @Param				request			body			user_experience.ExperienceIn		true		"project query params"
// @Success				201				{object}		utils.ApiResponses
// @Failure				400				{object}		utils.ApiResponses	"Invalid request"
// @Failure				401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure				403				{object}		utils.ApiResponses	"Forbidden"
// @Failure				500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router				/experience/{companyID}/create/{userID}	[post]
func (db Database) CreateUserExperience(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Parse and validate the company ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
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
	// Parse the incoming JSON request into a ExperienceIn struct
	experience := new(ExperienceIn)
	if err := ctx.ShouldBindJSON(experience); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new project in the database
	dbExperience := &domains.UserExperience{
		ID:                   uuid.New(),
		ProfessionalTraining: experience.ProfessionalTraining,
		UserID:               userID,
		CompanyID:            companyID,
	}
	// Créer le projet uniquement si le code de projet n'existe pas déjà dans la base de données
	if err := domains.Create(db.DB, dbExperience); err != nil {
		logrus.Error("Erreur lors de l'enregistrement des données dans la base de données. Erreur: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Répondre avec succès
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}


// ReadUserExperience	Handles the retrieval of all UserExperiences.
// @Summary        	Get UserExperience
// @Description    	Get all UserExperience.
// @Tags			UserExperience
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		path		string		true		"Company ID"
// @Success			200				{object}	user_experience.ExperiencePagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/experience/{companyID}	[get]
func (db Database) ReadUserExperience(ctx *gin.Context) {

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

	// Check if the user's value is among the allowed choices
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

	// Retrieve all experiences data from the database
	experiences, err := ReadAllPagination(db.DB, []domains.UserExperience{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all project data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.UserExperience{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a project structure as a response
	response := ExperiencePagination{}
	dataTableProject := []ExperienceTable{}
	for _, experience := range experiences {

		dataTableProject = append(dataTableProject, ExperienceTable{
			ID:                   experience.ID,
			ProfessionalTraining: experience.ProfessionalTraining,
		})
	}
	response.Items = dataTableProject
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadExperienceList	Handles the retrieval the list of all Experiences.
// @Summary        	Get list of  Experiences
// @Description    	Get list of all Experiences.
// @Tags			UserExperience
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Success			200					{array}			user_experience.ExperienceList
// @Failure			400			    	{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/experience/{companyID}/list	[get]
func (db Database) ReadExperiencesList(ctx *gin.Context) {

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

	// Retrieve all UserExperiences data from the database
	experiences, err := ReadAllList(db.DB, []domains.UserExperience{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all project data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a project structure as a response
	projectList := []ExperienceList{}
	for _, exp := range experiences {
		projectList = append(projectList, ExperienceList{
			ID: exp.ID,

			ProfessionalTraining: exp.ProfessionalTraining,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, projectList)
}

// ReadExperienceCount	Handles the retrieval the number of all experiences.
// @Summary        	Get number of  Experiences
// @Description    	Get number of all Experiences.
// @Tags			UserExperience
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true	"Company ID"
// @Success			200						{object}		user_experience.ExperienceCount
// @Failure			400						{object}		utils.ApiResponses	"Invalid request"
// @Failure			401						{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses	"Forbidden"
// @Failure			500						{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/experience/{companyID}/count	[get]
func (db Database) ReadUserExperiencesCount(ctx *gin.Context) {

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

	// Retrieve all project data from the database
	Usersexperiences, err := domains.ReadTotalCount(db.DB, &[]domains.UserExperience{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all project data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a project structure as a response
	ExperiencesCount := ExperienceCount{
		Count: Usersexperiences,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, ExperiencesCount)
}

// ReadUserExperience	Handles the retrieval of one UserExperience.
// @Summary        	Get UserExperience
// @Description    	Get one UserExperience.
// @Tags			 UserExperience
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"Experience ID"
// @Success			200					{object}		user_experience.ExperienceDetails
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/experience/{companyID}/{ID}/Get	[get]
func (db Database) ReadOneUserExperiences(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	ExperienceID, err := uuid.Parse(ctx.Param("ID"))
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

	// Retrieve project data from the database
	exp, err := 	ReadByID(db.DB, domains.UserExperience{}, ExperienceID)
	if err != nil {
		logrus.Error("Error retrieving project data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a project structure as a response
	details := ExperienceDetails{
		ID:                   exp.ID,
		ProfessionalTraining: exp.ProfessionalTraining,
		CompanyID:            exp.CompanyID,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateUserExperience		Handles the update of a  UserExperience .
// @Summary        	Update  UserExperience
// @Description    	Update  UserExperience .
// @Tags			 UserExperience
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true	"Company ID"
// @Param			ID					path			string				true	"experience  ID"
// @Param			request				body			user_experience.ExperienceIn		true	"Project query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/experience/{companyID}/{ID}/update	[put]
func (db Database) UpdateUserExperience(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the project ID from the request parameter
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

	// Parse the incoming JSON request into a ProjectIn struct
	experience := new(ExperienceIn)
	if err := ctx.ShouldBindJSON(experience); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the project with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.UserExperience{}, objectID); err != nil {
		logrus.Error("Error checking if the project with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Update the project data in the database
	dbExperience := &domains.UserExperience{

		ProfessionalTraining:  experience.ProfessionalTraining,

	}


	
	
	if err = domains.Update(db.DB, dbExperience, objectID); err != nil {
		logrus.Error("Error updating project data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())

}

// DeleteUserExperience	 	Handles the deletion of a UserExperience.
// @Summary        	Delete  UserExperience
// @Description    	Delete one UserExperience.
// @Tags			UserExperience
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"Experience ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/experience/{companyID}/{ID}	[delete]
func (db Database) DeleteExperience(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the project ID from the request parameter
	ExperienceID, err := uuid.Parse(ctx.Param("ID"))
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

	// Check if the user experience  with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.UserExperience{}, ExperienceID); err != nil {
		logrus.Error("Error checking if the project with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the experience  data from the database
	if err := domains.Delete(db.DB, &domains.UserExperience{}, ExperienceID); err != nil {
		logrus.Error("Error deleting project data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
