package exitPermission

import (
	"labs/constants"
	"labs/domains"
	"net/http"
	"strconv"
	"time"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateExitPermission 		Handles the creation of a new exitPermission.
// @Summary        	Create an exitPermission demande
// @Description    	Create a new exitPermission.
// @Tags			ExitPermission
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			request			body			exitPermission.ExitPermissionDemande		true		"ExitPermission query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/exitPermission	[post]
func (db Database) AddExitPermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	logrus.Error("userUd from session", session.UserID)

	// Parse the incoming JSON request into a ExitPermissionDemande struct
	exitPermission := new(ExitPermissionDemande)
	if err := ctx.ShouldBindJSON(exitPermission); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	layout := "2006-01-02"
	dt, err := time.Parse(layout, exitPermission.StartDate)
	if err != nil {
		logrus.Error("Error handling the date. Error: ", err.Error())
	}

	// Create a new exitPermission demande in the database
	dbPermission := &domains.ExitPermission{
		ID:          uuid.New(),
		ReleaseDate: time.Now(),
		Reason:      exitPermission.Reason,
		StartTime:   dt,
		ReturnTime:  exitPermission.ReturnDate,
		Type:        exitPermission.Type,
		UserID:      session.UserID,
	}
	if err := domains.Create(db.DB, dbPermission); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadExitPermission			Handles the retrieval of all exit permissions for a specific user .
// @Summary        		Get exit permission
// @Description    		Get all exit permission for a specific user .
// @Tags				ExitPermission
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID				path			string		true		"User ID"
// @Param			page			query		int					false	"Page"
// @Param			limit			query		int					false	"Limit"
// @Success				200					{array}			exitPermission.ExitPermissionPagination
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/exitPermission/{userID}	[get]
func (db Database) ReadAllExitPermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the exitPermission ID from the request parameter
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

	// Check if the ExitPermissionDemande's value is among the allowed choices
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

	// Retrieve all exitPermission data from the database
	exitPermission, err := ReadAllPagination(db.DB, []domains.ExitPermission{}, userID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all exitPermission data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.ExitPermission{}, "user_id", userID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a ExitPermission structure as a response
	response := ExitPermissionPagination{}
	listExitPermission := []ExitPermissionDetails{}
	for _, exitPermission := range exitPermission {

		listExitPermission = append(listExitPermission, ExitPermissionDetails{
			ID:        exitPermission.ID,
			Status:    exitPermission.Status,
			Reason:    exitPermission.Reason,
			UserID:    exitPermission.UserID,
			CreatedAt: exitPermission.CreatedAt,
		})
	}

	response.Items = listExitPermission
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)

}

//************************************************

// ReadExitPermissionCount	Handles the retrieval the number of all exitPermission.
// @Summary        			Get exitPermission count
// @Description    			Get all exitPermission count.
// @Tags					ExitPermission
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					userID				path			string		true		"User ID"
// @Success					200					{object}		exitPermission.ExitPermissionCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router					/exitPermission/user/{userID}/count	[get]
func (db Database) ReadExitPermissionCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the exitPermission ID from the request parameter
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

	// Retrieve all exitPermission data from the database
	exitPermission, err := domains.ReadTotalCount(db.DB, &domains.ExitPermission{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all exitPermission data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := ExitPermissionCount{
		Count: exitPermission,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}

// ReadExitPermission		Handles the retrieval of one exitPermission.
// @Summary        		Get exitPermission
// @Description    		Get one exitPermission.
// @Tags				ExitPermission
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string		true		"User ID"
// @Param				ID						path			string		true		"ExitPermission ID"
// @Success				200						{object}		exitPermission.ExitPermissionDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/exitPermission/user/{userID}/{ID}	[get]
func (db Database) ReadOneExitPermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the exitPermission ID from the request parameter
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

	// Retrieve exitPermission data from the database
	exitPermission, err := ReadByID(db.DB, domains.ExitPermission{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving exitPermission data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate an exitPermission structure as a response
	details := ExitPermissionDetails{
		ID:        exitPermission.ID,
		Reason:    exitPermission.Reason,
		Status:    exitPermission.Status,
		UserID:    exitPermission.UserID,
		CreatedAt: exitPermission.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateExitPermission 	Handles the update of a exitPermission.
// @Summary        		Update exitPermission
// @Description    		Update one exitPermission.
// @Tags				ExitPermission
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string							true		"User ID"
// @Param				ID						path			string							true		"ExitPermission ID"
// @Param				request					body			exitPermission.ExitPermissionIn	true		"ExitPermission query params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/exitPermission/user/{userID}/{ID}	[put]
func (db Database) UpdateExitPermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the exitPermission ID from the request parameter
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

	// Parse the incoming JSON request into a ExitPermission struct
	exitPermission := new(ExitPermissionIn)
	if err := ctx.ShouldBindJSON(exitPermission); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the exitPermission with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.ExitPermission{}, objectID); err != nil {
		logrus.Error("Error checking if the exitPermission with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the exitPermission data in the database
	dbExitPermission := &domains.ExitPermission{
		Status: exitPermission.Status,
	}
	if err = domains.Update(db.DB, dbExitPermission, objectID); err != nil {
		logrus.Error("Error updating exitPermission data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteExitPermission 	Handles the deletion of a exitPermission.
// @Summary        		Delete exitPermission
// @Description    		Delete one exitPermission.
// @Tags				ExitPermission
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string			true		"User ID"
// @Param				ID						path			string			true		"ExitPermission ID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"Internal Server Error"
// @Router				/exitPermission/user/{userID}/{ID}	[delete]
func (db Database) DeleteExitPermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the exitPermission ID from the request parameter
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

	// Delete the exitPermission data from the database
	if err := domains.Delete(db.DB, &domains.ExitPermission{}, objectID); err != nil {
		logrus.Error("Error deleting exitPermission data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
