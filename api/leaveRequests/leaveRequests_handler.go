package leaveRequests

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

// CreateLeaveRequest 		Handles the creation of a new LeaveRequests.
// @Summary        	Create LeaveRequests demande
// @Description    	Create a new LeaveRequests.
// @Tags			LeaveRequests
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			request			body			leaveRequests.LeaveRequestDemande		true		"LeaveRequests query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/LeaveRequests	[post]
func (db Database) AddLeave(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	logrus.Error("userUd from session", session.UserID)

	// Parse the incoming JSON request into a UserIn struct
	LeaveRequests := new(LeaveRequestDemande)
	if err := ctx.ShouldBindJSON(LeaveRequests); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new LeaveRequests demande in the database
	dbLeave := &domains.LeaveRequests{
		ID:        uuid.New(),
		StartDate: LeaveRequests.StartDate,
		EndDate:   LeaveRequests.EndDate,
		Type:      LeaveRequests.Type,
		Reason:    LeaveRequests.Reason,
		UserID:    session.UserID,
	}
	if err := domains.Create(db.DB, dbLeave); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadLeaveRequests			Handles the retrieval of all LeaveRequests for a specific user .
// @Summary        		Get LeaveRequests
// @Description    		Get all LeaveRequests.
// @Tags				LeaveRequests
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID				path			string		true		"User ID"
// @Param				page			query		int					false       "Page"
// @Param				limit			query		int					false	    "Limit"
// @Success				200					{array}			leaveRequests.LeaveRequestDetails
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/LeaveRequests/{userID}	[get]
func (db Database) ReadLeave(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the LeaveRequests ID from the request parameter
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

	// Check if the LeaveRequest's value is among the allowed choices
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
	LeaveRequests, err := ReadAllPagination(db.DB, []domains.LeaveRequests{}, userID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.LeaveRequests{}, "user_id", userID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a LeaveRequests structure as a response
	response := LeavePagination{}
	listLeave := []LeaveRequestDetails{}
	for _, LeaveRequests := range LeaveRequests {

		listLeave = append(listLeave, LeaveRequestDetails{
			ID:        LeaveRequests.ID,
			StartDate: LeaveRequests.StartDate,
			EndDate:   LeaveRequests.EndDate,
			Status:    LeaveRequests.Status,
			Reason:    LeaveRequests.Reason,
			CreatedAt: LeaveRequests.CreatedAt,
		})
	}

	response.Items = listLeave
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)

}

// ReadLeaveRequestsCount	Handles the retrieval the number of all LeaveRequests.
// @Summary        			Get LeaveRequests count
// @Description    			Get all LeaveRequests count.
// @Tags					LeaveRequests
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					userID				path			string		true		"User ID"
// @Success					200					{object}		leaveRequests.LeaveRequestCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router					/LeaveRequests/{userID}/count	[get]
func (db Database) ReadLeaveCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the LeaveRequests ID from the request parameter
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

	// Retrieve all LeaveRequests data from the database
	LeaveRequests, err := domains.ReadTotalCount(db.DB, &domains.LeaveRequests{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := LeaveRequestCount{
		Count: LeaveRequests,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}

// ReadLeaveRequests		Handles the retrieval of one LeaveRequests.
// @Summary        		Get LeaveRequests
// @Description    		Get one LeaveRequests.
// @Tags				LeaveRequests
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string		true		"User ID"
// @Param				ID						path			string		true		"LeaveRequests ID"
// @Success				200						{object}		leaveRequests.LeaveRequestDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/LeaveRequests/{userID}/{ID}	[get]
func (db Database) ReadOneLeave(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the LeaveRequests ID from the request parameter
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

	// Retrieve LeaveRequests data from the database
	role, err := ReadByID(db.DB, domains.LeaveRequests{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving LeaveRequests data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a LeaveRequest structure as a response
	details := LeaveRequestDetails{
		ID:        role.ID,
		StartDate: role.StartDate,
		EndDate:   role.EndDate,
		Status:    role.Status,
		Reason:    role.Reason,
		CreatedAt: role.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateLeaveRequest 	Handles the update of a LeaveRequests.
// @Summary        		Update LeaveRequests
// @Description    		Update one LeaveRequests.
// @Tags				LeaveRequests
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string							true		"User ID"
// @Param				ID						path			string							true		"LeaveRequests ID"
// @Param				request					body			leaveRequests.LeaveRequestIn	true		"LeaveRequests query params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/LeaveRequests/{userID}/{ID}	[put]
func (db Database) UpdateLeave(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the LeaveRequests ID from the request parameter
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

	// Parse the incoming JSON request into a  LeaveIn struct
	LeaveRequests := new(LeaveRequestIn)
	if err := ctx.ShouldBindJSON(LeaveRequests); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the LeaveRequests with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.LeaveRequests{}, objectID); err != nil {
		logrus.Error("Error checking if the LeaveRequests with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the LeaveRequests data in the database
	dbLeave := &domains.LeaveRequests{
		Status: LeaveRequests.Status,
	}
	if err = domains.Update(db.DB, dbLeave, objectID); err != nil {
		logrus.Error("Error updating LeaveRequests data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteLeaveRequest 	Handles the deletion of a LeaveRequests.
// @Summary        		Delete LeaveRequests
// @Description    		Delete one LeaveRequests.
// @Tags				LeaveRequests
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string			true		"User ID"
// @Param				ID						path			string			true		"LeaveRequests ID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"Internal Server Error"
// @Router				/LeaveRequests/{userID}/{ID}	[delete]
func (db Database) DeleteLeave(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the LeaveRequests ID from the request parameter
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

	// Delete the LeaveRequest data from the database
	if err := domains.Delete(db.DB, &domains.LeaveRequests{}, objectID); err != nil {
		logrus.Error("Error deleting LeaveRequest data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
