package companies

import (
	"labs/constants"
	"labs/domains"
	"log"
	"net/http"
	"strconv"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateCompany 	Handles the creation of a new company.
// @Summary        	Create company
// @Description    	Create a new company.
// @Tags			Companies
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			request			body			companies.CompanyIn	true	"Company query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/companies		[post]
func (db Database) CreateCompany(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse the incoming JSON request into a CompanyIn struct
	company := new(CompanyIn)
	if err := ctx.ShouldBindJSON(company); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new user in the database
	dbCompany := &domains.Companies{
		ID:              uuid.New(),
		Name:            company.Name,
		CreatedByUserID: session.UserID,
	}
	if err := domains.Create(db.DB, dbCompany); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// ReadCompanies 	Handles the retrieval of all companies.
// @Summary        	Get companies
// @Description    	Get all companies.
// @Tags			Companies
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int					false	"Page"
// @Param			limit			query		int					false	"Limit"
// @Success			200				{object}	companies.CompaniesPagination
// @Failure			400				{object}	utils.ApiResponses			"Invalid request"
// @Failure			401				{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses			"Forbidden"
// @Failure			500				{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/companies		[get]
func (db Database) ReadCompanies(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	log.Println("session")
	log.Println(session)
	log.Println("session")

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

	// Retrieve all company data from the database
	companies, err := ReadAllPagination(db.DB, []domains.Companies{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Companies{}, "id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a company structure as a response
	response := CompaniesPagination{}
	listCompany := []CompaniesTable{}
	for _, company := range companies {

		listCompany = append(listCompany, CompaniesTable{
			ID:        company.ID,
			Name:      company.Name,
			Email:     company.Email,
			CreatedAt: company.CreatedAt,
		})
	}
	response.Items = listCompany
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadCompany 		Handles the retrieval of one company.
// @Summary        	Get company
// @Description    	Get one company.
// @Tags			Companies
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			ID   			path      	string		true		"Company ID"
// @Success			200				{object}	companies.CompaniesDetails
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/companies/{ID}	[get]
func (db Database) ReadCompany(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, objectID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the company data by ID from the database
	company, err := ReadByID(db.DB, domains.Companies{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a company structure as a response
	details := CompaniesDetails{
		ID:        company.ID,
		Name:      company.Name,
		Email:     company.Email,
		Website:   company.Website,
		CreatedAt: company.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateCompany 	Handles the update of a company.
// @Summary        	Update company
// @Description    	Update company.
// @Tags			Companies
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			ID   			path      		string						true	"Company ID"
// @Param			request			body			companies.CompanyIn		true	"Company query params"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses				"Invalid request"
// @Failure			401				{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses				"Forbidden"
// @Failure			500				{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/companies/{ID}	[put]
func (db Database) UpdateCompany(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, objectID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a CompanyIn struct
	company := new(CompanyIn)
	if err := ctx.ShouldBindJSON(company); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the company with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Companies{}, objectID); err != nil {
		logrus.Error("Error checking if the company with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the company data in the database
	dbCompany := &domains.Companies{
		Name: company.Name,
	}
	if err := domains.Update(db.DB, dbCompany, objectID); err != nil {
		logrus.Error("Error updating company data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteCompany 	Handles the deletion of a company.
// @Summary        	Delete company
// @Description    	Delete one company.
// @Tags			Companies
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			ID   			path      		string		true			"Company ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/companies/{ID}	[delete]
func (db Database) DeleteCompany(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, objectID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the company with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Companies{}, objectID); err != nil {
		logrus.Error("Error checking if the company with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the company data from the database
	if err := domains.Delete(db.DB, &domains.Companies{}, objectID); err != nil {
		logrus.Error("Error deleting company data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
