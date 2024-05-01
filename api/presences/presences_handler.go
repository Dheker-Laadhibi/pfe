package presences

import (
	"labs/constants"
	"labs/domains"
	"net/http"
	"strconv"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//create presence

// CreatePresence 	Handles the creation of a new presence.
// @Summary        	Create presence
// @Description    	Create a new presence.
// @Tags			Presences
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID	    path			string				true	"companyID"
// @Param			userID	        path			string				true	"userID"
// @Param			request			body			presences.PresencesIn	true "Presence query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Presenceal Server Error"
// @Router			/presences/{companyID}/{userID}	[post]
func (db Database) CreatePresence(ctx *gin.Context) {

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
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified user
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a PresenceIn struct
	presence := new(PresencesIn)
	if err := ctx.ShouldBindJSON(presence); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new intern in the database
	dbPresence := &domains.Presences{
		ID:        uuid.New(),
		Matricule: presence.Matricule,
		Check:     presence.Check,
		UserID: userID, //IDUser

	}
	if err := domains.Create(db.DB, dbPresence); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}



// ReadPresences   Handles the retrieval of all presences.
// @Summary        Get presences
// @Description    Get all presences.
// @Tags		   Presences
// @Produce		   json
// @Security 	   ApiKeyAuth
// @Param		   page			   query		int			false		"Page"
// @Param		   limit		   query		int			false		"Limit"
// @Param		   userID		   path			string		true		"User ID"
// @Success		   200			   {array}		presences.PresencesDetails
// @Failure		   400			   {object}		utils.ApiResponses		"Invalid request"
// @Failure		   401			   {object}		utils.ApiResponses		"Unauthorized"
// @Failure		   403			   {object}		utils.ApiResponses		"Forbidden"
// @Failure		   500			   {object}		utils.ApiResponses		"Presenceal Server Error"
// @Router		   /presences/All/{userID}	[get]
func (db Database) ReadPresences(ctx *gin.Context) {

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

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
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
	// CheckEmployeeSession checks if the user's session matches the specified user and company.
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
// Retrieve all user data from the database
presences, err := ReadAllPagination(db.DB, []domains.Presences{}, session.UserID, limit, offset)
if err != nil {
	logrus.Error("Error occurred while finding all user data. Error: ", err)
	utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
	return
}
count, err := domains.ReadTotalCount(db.DB, &domains.Presences{}, "user_id", session.UserID)
if err != nil {
	logrus.Error("Error occurred while finding total count. Error: ", err)
	utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
	return
}

	
// Generate the presence's   structure as a response
response := PresencesPagination{}
dataTablePresence := []PresencesTable{}
for _, presence := range presences {

	dataTablePresence = append(dataTablePresence, PresencesTable{
		ID:        presence.ID,
		Check:presence.Check ,
		Matricule:presence.Matricule,
		
		
	})
}
response.Items = dataTablePresence
response.Page = uint(page)
response.Limit = uint(limit)
response.TotalCount = count


	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}


// ReadPresencesCount	Handles the retrieval the number of all presences.
// @Summary        		Get presences count
// @Description    		Get all presences count.
// @Tags				Presences
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID				path			string		true		"User ID"
// @Success				200					{object}		presences.PresencesCount
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/presences/count/{userID}	[get] 
func (db Database) ReadPresencesCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all presence data from the database
	presences, err := domains.ReadTotalCount(db.DB, &domains.Presences{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all presence data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := PresencesCount{
		Count: presences,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)




}



// ReadPresence		    Handles the retrieval of one presence.
// @Summary        		Get presence
// @Description    		Get one presence.
// @Tags				Presences
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string		true		"User ID"
// @Param				ID						path			string		true		"Presence ID"
// @Success				200						{object}		presences.PresencesDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/presences/get/{ID}/{userID}	[get]
func (db Database) ReadPresence(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the presence ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve presence data from the database
	presence, err := ReadByID(db.DB, domains.Presences{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving presence data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a   presence   structure as a response
	details := PresencesDetails{
		ID:        presence.ID,
		Matricule: presence.Matricule,

		Check: presence.Check,
		UserID: userID,

		
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}
