package candidats

import (
	"labs/constants"
	"labs/domains"
	"net/http"
	"strconv"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Createcandidate 		Handles the creation of a new candidate.
// @Summary        	Create candidate
// @Description    	Create a new candidate.
// @Tags			Condidats
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				        true		"Company ID"
// @Param			request			body			condidats.CondidatIn		true		"condidat query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/condidats/{companyID}	[post]
func (db Database) Createcandidate(ctx *gin.Context) {

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

	// Parse the incoming JSON request into a candidatIn struct
	condidat := new(CondidatIn)
	if err := ctx.ShouldBindJSON(condidat); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Hash the user's password
	hash, _ := bcrypt.GenerateFromPassword([]byte(condidat.Password), bcrypt.DefaultCost)
	// Create a new candidate in the database
	dbCondidat := &domains.Condidats{
		ID:             uuid.New(),
		Firstname:      condidat.Firstname,
		Lastname:       condidat.Lastname,
		Email:          condidat.Email,
		Password:       string(hash),
		Adress:         condidat.Adress,
		University:     condidat.University,
		Educationlevel: condidat.Educationlevel,
		CompanyID:      condidat.CompanyID,
	}
	if err := domains.Create(db.DB, dbCondidat); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadCandidats 	Handles the retrieval of all Candidats .
// @Summary        	Get Candidats
// @Description    	Get all Candidats .
// @Tags			Condidats
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page				query		int				false		"Page"
// @Param			limit				query		int				false		"Limit"
// @Param			companyID			path		string			true		"Company ID"
// @Success			200					{object}	condidats.CondidtasPagination
// @Failure			400					{object}	utils.ApiResponses			"Invalid request"
// @Failure			401					{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}	utils.ApiResponses			"Forbidden"
// @Failure			500					{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/condidats/{companyID}	[get]
func (db Database) ReadCandidats(ctx *gin.Context) {

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

	// Retrieve all Candidats  data from the database
	condidats, err := ReadAllPagination(db.DB, []domains.Condidats{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Condidats{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a user structure as a response
	response := CondidtasPagination{}
	dataTableCondidat := []CondidatsTable{}
	for _, condidat := range condidats {

		dataTableCondidat = append(dataTableCondidat, CondidatsTable{
			ID:        condidat.ID,
			Firstname: condidat.Firstname,
			Lastname:  condidat.Lastname,
			Email:     condidat.Email,
		})
	}
	response.Items = dataTableCondidat
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadCandidatsList	Handles the retrieval the list of all Candidats.
// @Summary        	Get list of  Candidats
// @Description    	Get list of all Candidats.
// @Tags			Condidats
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Success			200						{array}			condidats.CondidatsList
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/condidats/{companyID}/list	[get]
func (db Database) ReadCondidatsList(ctx *gin.Context) {

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

	// Retrieve all condidat data from the database
	condidats, err := ReadAllList(db.DB, []domains.Condidats{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a Candidats  structure as a response
	condidatsList := []CondidatsList{}
	for _, condidat := range condidats {

		condidatsList = append(condidatsList, CondidatsList{
			ID:        condidat.ID,
			Firstname: condidat.Firstname,
			Lastname:  condidat.Lastname,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, condidatsList)
}

// ReadCandidatsCount	Handles the retrieval the number of all Candidats .
// @Summary        	Get number of  Candidats
// @Description    	Get number of all Candidats .
// @Tags			Condidats
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true		"Company ID"
// @Success			200						{object}		condidats.CondidatsCount
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/condidats/{companyID}/count	[get]
func (db Database) ReadCandidatsCount(ctx *gin.Context) {

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

	// Retrieve all condidat data from the database
	condidats, err := domains.ReadTotalCount(db.DB, &[]domains.Condidats{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a condidat structure as a response
	condidatsCount := CondidatsCount{
		Count: condidats,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, condidatsCount)
}

// Readcandidat		Handles the retrieval of one candidat.
// @Summary        	Get candidat
// @Description    	Get one candidat.
// @Tags			Condidats
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Param			ID						path			string			true	"candidat ID"
// @Success			200						{object}		condidats.CondidatDetails
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/condidats/{companyID}/{ID}	[get]
func (db Database) Readcandidat(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the condidats ID from the request parameter
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

	// Retrieve candidat data from the database
	condidat, err := ReadByID(db.DB, domains.Condidats{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving condidat data from the database. Error: ", err.Error())
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

	// Generate a Candidat structure as a response
	details := CondidatDetails{
		ID:               condidat.ID,
		Firstname:        condidat.Firstname,
		Lastname:         condidat.Lastname,
		LevelOfEducation: condidat.Educationlevel,
		University:       condidat.University,
		CompanyID:        condidat.CompanyID,
		CompanyName:      companyName,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateCandidt	Handles the update of a candidat.
// @Summary        	Update candidat
// @Description    	Update one candidat.
// @Tags			Condidats
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true		"Company ID"
// @Param			ID					path			string				true		"candidat ID"
// @Param			request				body			condidats.CondidatIn		true		"candidat query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses				"Invalid request"
// @Failure			401					{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses				"Forbidden"
// @Failure			500					{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/condidats/{companyID}/{ID}	[put]
func (db Database) Updatecandidat(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the condidat ID from the request parameter
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

	// Parse the incoming JSON request into a candidatIn struct
	condidat := new(CondidatIn)
	if err := ctx.ShouldBindJSON(condidat); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the candidat with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Condidats{}, objectID); err != nil {
		logrus.Error("Error checking if the condidats with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the candidat data in the database
	dbCondidat := &domains.Condidats{
		Firstname:      condidat.Firstname,
		Lastname:       condidat.Lastname,
		Email:          condidat.Email,
		Adress:         condidat.Adress,
		University:     condidat.University,
		Educationlevel: condidat.Educationlevel,
	}
	if err = domains.Update(db.DB, dbCondidat, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// Deletecandidat 	Handles the deletion of a candidat	.
// @Summary        	Delete candidat
// @Description    	Delete one candidat	.
// @Tags			Condidats
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true		"Company ID"
// @Param			ID					path			string			true		"candidat ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/condidats/{companyID}/{ID}	[delete]
func (db Database) DeleteCondidat(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the Condidats ID from the request parameter
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

	// Check if the candidat with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Condidats{}, objectID); err != nil {
		logrus.Error("Error checking if the condidat with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the candidat data from the database
	if err := domains.Delete(db.DB, &domains.Condidats{}, objectID); err != nil {
		logrus.Error("Error deleting condidat data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// SigninCandidat 	Handles the candidat signin process.
// @Summary			Signin
// @Description		Authenticate and log in a uscandidater.
// @Tags			Condidats
// @Accept			json
// @Produce			json
// @Param			companyID			path			string			true		"Company ID"
// @Param			request		body		condidats.Signin		true	"candidat query params"
// @Success			200			{object}	condidats.LoggedIn
// @Failure			400			{object}	utils.ApiResponses		"Invalid request"
// @Failure			401			{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403			{object}	utils.ApiResponses		"Forbidden"
// @Failure			500			{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/condidats/{companyID}/Signin	[post]
func (db Database) SigninCandidat(ctx *gin.Context) {

	// Parse the incoming JSON request into a Signin struct
	candidat := new(Signin)
	if err := ctx.ShouldBindJSON(candidat); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve user data by email
	data, err := ReadByEmailActive(db.DB, candidat.Email)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Compare the entered password with the stored password
	if isTrue := utils.ComparePassword(data.Password, candidat.Password); !isTrue {
		logrus.Error("Password comparison failed.")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNAUTHORIZED, utils.Null())
		return
	}

	// Prepare the response with candidat details
	response := LoggedIn{
		ID:        data.ID,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		CompanyID: data.CompanyID,
	}
	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}
