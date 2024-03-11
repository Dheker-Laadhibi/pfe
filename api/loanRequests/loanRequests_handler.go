package loanRequests

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

// CreateLoanRequest 		Handles the creation of a new loanRequests.
// @Summary        	Create an loanRequests demande
// @Description    	Create a new loanRequests.
// @Tags			LoanRequests
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			request			body			loanRequests.LoanRequestDemande		true		"LoanRequests query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/loanRequests	[post]
func (db Database) AddLoanRequests(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	logrus.Error("userUd from session", session.UserID)

	// Parse the incoming JSON request into a LoanRequestDemande struct
	loanRequests := new(LoanRequestDemande)
	if err := ctx.ShouldBindJSON(loanRequests); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new loanRequests demande in the database
	dbLoanRequests := &domains.LoanRequests{
		ID:            uuid.New(),
		LoanAmount:    loanRequests.LoanAmount,
		LoanDuration:  loanRequests.LoanDuration,
		InterestRate:  loanRequests.InterestRate,
		ReasonForLoan: loanRequests.ReasonForLoan,
		PathDocument:  loanRequests.PathDocument,
		UserID:        session.UserID,
	}
	if err := domains.Create(db.DB, dbLoanRequests); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadLoanRequest			Handles the retrieval of all loanRequests for a specific user .
// @Summary        			Get LoanRequests
// @Description    			Get all LoanRequest for a specific user .
// @Tags					LoanRequests
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					userID				path			string		true		"User ID"
// @Param					page			query		int					false	"Page"
// @Param					limit			query		int					false	"Limit"
// @Success				200					{array}			loanRequests.LoanRequestPagination
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/loanRequests/{userID}	[get]
func (db Database) ReadAllLoanRequests(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the loanRequests ID from the request parameter
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

	// Check if the loanRequest's value is among the allowed choices
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

	// Retrieve all loanRequests data from the database
	loanRequests, err := ReadAllPagination(db.DB, []domains.LoanRequests{}, userID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all loanRequests data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.LoanRequests{}, "user_id", userID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a LoanRequests structure as a response
	response := LoanRequestPagination{}
	listLoanRequest := []LoanRequestDetails{}
	for _, loanRequests := range loanRequests {

		listLoanRequest = append(listLoanRequest, LoanRequestDetails{
			ID:            loanRequests.ID,
			LoanAmount:    loanRequests.LoanAmount,
			LoanDuration:  loanRequests.LoanDuration,
			InterestRate:  loanRequests.InterestRate,
			ReasonForLoan: loanRequests.ReasonForLoan,
			PathDocument:  loanRequests.PathDocument,
			Status:        loanRequests.Status,
			UserID:        loanRequests.UserID,
			CreatedAt:     loanRequests.CreatedAt,
		})
	}

	response.Items = listLoanRequest
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)

}

//************************************************

// ReadLoanRequestCount	Handles the retrieval the number of all loanRequests.
// @Summary        			Get loanRequests count
// @Description    			Get all loanRequests count.
// @Tags					LoanRequests
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					userID				path			string		true		"User ID"
// @Success					200					{object}		loanRequests.LoanRequestCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router					/loanRequests/user/{userID}/count	[get]
func (db Database) ReadLoanRequestsCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the loanRequests ID from the request parameter
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

	// Retrieve all loanRequests data from the database
	loanRequests, err := domains.ReadTotalCount(db.DB, &domains.LoanRequests{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all loanRequests data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := LoanRequestCount{
		Count: loanRequests,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}

// ReadLoanRequest		Handles the retrieval of one loanRequests.
// @Summary        		Get loanRequests
// @Description    		Get one loanRequests.
// @Tags				LoanRequests
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string		true		"User ID"
// @Param				ID						path			string		true		"LoanRequests ID"
// @Success				200						{object}		loanRequests.LoanRequestDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/loanRequests/user/{userID}/{ID}	[get]
func (db Database) ReadOneLoanRequests(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the loanRequests ID from the request parameter
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

	// Retrieve loanRequests data from the database
	loanRequests, err := ReadByID(db.DB, domains.LoanRequests{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving loanRequests data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a loanRequest structure as a response
	details := LoanRequestDetails{
		ID:            loanRequests.ID,
		LoanAmount:    loanRequests.LoanAmount,
		LoanDuration:  loanRequests.LoanDuration,
		InterestRate:  loanRequests.InterestRate,
		ReasonForLoan: loanRequests.ReasonForLoan,
		PathDocument:  loanRequests.PathDocument,
		Status:        loanRequests.Status,
		UserID:        loanRequests.UserID,
		CreatedAt:     loanRequests.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateLoanRequest 	Handles the update of a loanRequests.
// @Summary        		Update loanRequests
// @Description    		Update one loanRequests.
// @Tags				LoanRequests
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string							true		"User ID"
// @Param				ID						path			string							true		"LoanRequests ID"
// @Param				request					body			loanRequests.LoanRequestIn	true		"LoanRequests query params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/loanRequests/user/{userID}/{ID}	[put]
func (db Database) UpdateLoanRequests(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the loanRequests ID from the request parameter
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

	// Parse the incoming JSON request into a LoanRequest struct
	loanRequests := new(LoanRequestIn)
	if err := ctx.ShouldBindJSON(loanRequests); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the loanRequests with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.LoanRequests{}, objectID); err != nil {
		logrus.Error("Error checking if the loanRequests with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the loanRequests data in the database
	dbExitPermission := &domains.LoanRequests{
		Status: loanRequests.Status,
	}
	if err = domains.Update(db.DB, dbExitPermission, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteExitPermission 	Handles the deletion of a loanRequests.
// @Summary        		Delete loanRequests
// @Description    		Delete one loanRequests.
// @Tags				LoanRequests
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string			true		"User ID"
// @Param				ID						path			string			true		"LoanRequests ID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"Internal Server Error"
// @Router				/loanRequests/user/{userID}/{ID}	[delete]
func (db Database) DeleteLoanRequests(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the loanRequests ID from the request parameter
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

	// Delete the loanRequest data from the database
	if err := domains.Delete(db.DB, &domains.LoanRequests{}, objectID); err != nil {
		logrus.Error("Error deleting loanRequest data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
