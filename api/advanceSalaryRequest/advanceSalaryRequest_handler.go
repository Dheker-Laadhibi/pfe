package advanceSalaryRequest

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

// CreateAdvanceSalaryRequests		Handles the creation of a new AdvanceSalaryRequests.
// @Summary        	Create AdvanceSalaryRequestsdemande
// @Description    	Create a new AdvanceSalaryRequests.
// @Tags			AdvanceSalaryRequests
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			request			body			advanceSalaryRequest.AdvanceSalaryRequestDemande		true		"AdvanceSalaryRequestsquery params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/AdvanceSalaryRequests	[post]
func (db Database) AddAdvanceSalaryRequest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	logrus.Error("userUd from session", session.UserID)

	// Parse the incoming JSON request into a AdvanceSalaryRequestDemande struct
	advanceSalaryRequests := new(AdvanceSalaryRequestDemande)
	if err := ctx.ShouldBindJSON(advanceSalaryRequests); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new AdvanceSalaryRequests demande in the database
	dbAdvanceSalary := &domains.AdvanceSalaryRequests{
		ID:     uuid.New(),
		Amount: advanceSalaryRequests.Amount,
		Reason: advanceSalaryRequests.Reason,
		UserID: session.UserID,
	}
	if err := domains.Create(db.DB, dbAdvanceSalary); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadAdvanceSalaryRequests			Handles the retrieval of all AdvanceSalaryRequests for a specific user .
// @Summary        		Get AdvanceSalaryRequests
// @Description    		Get all AdvanceSalaryRequests.
// @Tags				AdvanceSalaryRequests
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID				path			string		true		"User ID"
// @Param				page			query		int					false       "Page"
// @Param				limit			query		int					false	    "Limit"
// @Success				200					{array}			advanceSalaryRequest.AdvanceSalaryRequestDetails
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/AdvanceSalaryRequests/{userID}	[get]
func (db Database) ReadAdvanceSalaryRequest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the AdvanceSalaryRequestsID from the request parameter
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

	// Check if the advanceSalaryRequest's value is among the allowed choices
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

	// Retrieve all company data from the database
	AdvanceSalaryRequests, err := ReadAllPagination(db.DB, []domains.AdvanceSalaryRequests{}, userID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.AdvanceSalaryRequests{}, "user_id", userID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a AdvanceSalaryRequest structure as a response
	response := AdvanceSalaryRequestPagination{}
	listAdvanceSalaryRequest := []AdvanceSalaryRequestDetails{}
	for _, advanceSalaryRequests := range AdvanceSalaryRequests {

		listAdvanceSalaryRequest = append(listAdvanceSalaryRequest, AdvanceSalaryRequestDetails{
			ID:        advanceSalaryRequests.ID,
			Amount:    advanceSalaryRequests.Amount,
			Status:    advanceSalaryRequests.Status,
			Reason:    advanceSalaryRequests.Reason,
			CreatedAt: advanceSalaryRequests.CreatedAt,
		})
	}

	response.Items = listAdvanceSalaryRequest
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)

}

// ReadAdvanceSalaryRequestsCount	Handles the retrieval the number of all AdvanceSalaryRequests.
// @Summary        			Get AdvanceSalaryRequestscount
// @Description    			Get all AdvanceSalaryRequestscount.
// @Tags					AdvanceSalaryRequests
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					userID				path			string		true		"User ID"
// @Success					200					{object}		advanceSalaryRequest.AdvanceSalaryRequestCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router					/AdvanceSalaryRequests/{userID}/count	[get]
func (db Database) ReadAdvanceSalaryRequestCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the AdvanceSalaryRequestsID from the request parameter
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

	// Retrieve all AdvanceSalaryRequest data from the database
	AdvanceSalaryRequests, err := domains.ReadTotalCount(db.DB, &domains.AdvanceSalaryRequests{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := AdvanceSalaryRequestCount{
		Count: AdvanceSalaryRequests,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}

// ReadAdvanceSalaryRequests		Handles the retrieval of one AdvanceSalaryRequests.
// @Summary        		Get AdvanceSalaryRequests
// @Description    		Get one AdvanceSalaryRequests.
// @Tags				AdvanceSalaryRequests
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string		true		"User ID"
// @Param				ID						path			string		true		"AdvanceSalaryRequestsID"
// @Success				200						{object}		advanceSalaryRequest.AdvanceSalaryRequestDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/AdvanceSalaryRequests/{userID}/{ID}	[get]
func (db Database) ReadOneAdvanceSalaryRequest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the AdvanceSalaryRequestsID from the request parameter
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

	// Retrieve AdvanceSalaryRequest data from the database
	advanceSalaryRequest, err := ReadByID(db.DB, domains.AdvanceSalaryRequests{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving AdvanceSalaryRequest data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a user structure as a response
	details := AdvanceSalaryRequestDetails{
		ID:        advanceSalaryRequest.ID,
		Amount:    advanceSalaryRequest.Amount,
		Status:    advanceSalaryRequest.Status,
		Reason:    advanceSalaryRequest.Reason,
		CreatedAt: advanceSalaryRequest.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateAdvanceSalaryRequest 	Handles the update of a AdvanceSalaryRequests.
// @Summary        		Update AdvanceSalaryRequests
// @Description    		Update one AdvanceSalaryRequests.
// @Tags				AdvanceSalaryRequests
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string							true		"User ID"
// @Param				ID						path			string							true		"AdvanceSalaryRequestsID"
// @Param				request					body			advanceSalaryRequest.AdvanceSalaryRequestIn	true		"AdvanceSalaryRequestsquery params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/AdvanceSalaryRequests/{userID}/{ID}	[put]
func (db Database) UpdateAdvanceSalaryRequest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the AdvanceSalaryRequestsID from the request parameter
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

	// Parse the incoming JSON request into a  advanceSalaryRequestIn struct
	AdvanceSalaryRequests := new(AdvanceSalaryRequestIn)
	if err := ctx.ShouldBindJSON(AdvanceSalaryRequests); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the AdvanceSalaryRequestswith the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.AdvanceSalaryRequests{}, objectID); err != nil {
		logrus.Error("Error checking if the AdvanceSalaryRequestswith the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the AdvanceSalaryRequestsdata in the database
	dbAdvanceSalaryRequest := &domains.AdvanceSalaryRequests{
		Status: AdvanceSalaryRequests.Status,
	}
	if err = domains.Update(db.DB, dbAdvanceSalaryRequest, objectID); err != nil {
		logrus.Error("Error updating AdvanceSalaryRequestsdata in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteAdvanceSalaryRequest 	Handles the deletion of a AdvanceSalaryRequests.
// @Summary        		Delete AdvanceSalaryRequests
// @Description    		Delete one LAdvanceSalaryRequests.
// @Tags				AdvanceSalaryRequests
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string			true		"User ID"
// @Param				ID						path			string			true		"AdvanceSalaryRequestsID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"Internal Server Error"
// @Router				/AdvanceSalaryRequests/{userID}/{ID}	[delete]
func (db Database) DeleteAdvanceSalaryRequest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the AdvanceSalaryRequestsID from the request parameter
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

	// Delete the user data from the database
	if err := domains.Delete(db.DB, &domains.AdvanceSalaryRequests{}, objectID); err != nil {
		logrus.Error("Error deleting AdvanceSalaryRequestsdata from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
