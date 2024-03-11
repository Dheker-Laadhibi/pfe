package roles

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

// CreateRole 		Handles the creation of a new role.
// @Summary        	Create role
// @Description    	Create a new role.
// @Tags			Roles
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				true		"Company ID"
// @Param			request			body			roles.RolesIn		true		"Role query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/roles/{companyID}	[post]
func (db Database) CreateRole(ctx *gin.Context) {

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

	// Parse the incoming JSON request into a RoleIn struct
	role := new(RolesIn)
	if err := ctx.ShouldBindJSON(role); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new role in the database
	dbRole := &domains.Roles{
		ID:              uuid.New(),
		Name:            role.Name,
		OwningCompanyID: companyID,
		CreatedByUserID: session.UserID,
	}
	if err := domains.Create(db.DB, dbRole); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadRoles 		Handles the retrieval of all roles.
// @Summary        	Get roles
// @Description    	Get all roles.
// @Tags			Roles
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page				query		int				false		"Page"
// @Param			limit				query		int				false		"Limit"
// @Param			companyID			path		string			true		"Company ID"
// @Success			200					{object}	roles.RolesPagination
// @Failure			400					{object}	utils.ApiResponses			"Invalid request"
// @Failure			401					{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}	utils.ApiResponses			"Forbidden"
// @Failure			500					{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/roles/{companyID}	[get]
func (db Database) ReadRoles(ctx *gin.Context) {

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

	// Retrieve all role data from the database
	roles, err := ReadAllPagination(db.DB, []domains.Roles{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Roles{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a user structure as a response
	response := RolesPagination{}
	dataTableRole := []RolesTable{}
	for _, role := range roles {

		dataTableRole = append(dataTableRole, RolesTable{
			ID:        role.ID,
			Name:      role.Name,
			CreatedAt: role.CreatedAt,
		})
	}
	response.Items = dataTableRole
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadRolesList	Handles the retrieval the list of all roles.
// @Summary        	Get list of  roles
// @Description    	Get list of all roles.
// @Tags			Roles
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Success			200						{array}			roles.RolesList
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/roles/{companyID}/list	[get]
func (db Database) ReadRolesList(ctx *gin.Context) {

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

	// Retrieve all role data from the database
	roles, err := ReadAllList(db.DB, []domains.Roles{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a role structure as a response
	rolesList := []RolesList{}
	for _, role := range roles {

		rolesList = append(rolesList, RolesList{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, rolesList)
}

// ReadRolesCount	Handles the retrieval the number of all roles.
// @Summary        	Get number of  roles
// @Description    	Get number of all roles.
// @Tags			Roles
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true		"Company ID"
// @Success			200						{object}		roles.RolesCount
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/roles/{companyID}/count	[get]
func (db Database) ReadRolesCount(ctx *gin.Context) {

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

	// Retrieve all role data from the database
	roles, err := domains.ReadTotalCount(db.DB, &[]domains.Roles{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a role structure as a response
	rolesCount := RolesCount{
		Count: roles,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, rolesCount)
}

// ReadRole 		Handles the retrieval of one role.
// @Summary        	Get role
// @Description    	Get one role.
// @Tags			Roles
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Param			ID						path			string			true	"Role ID"
// @Success			200						{object}		roles.RolesDetails
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/roles/{companyID}/{ID}	[get]
func (db Database) ReadRole(ctx *gin.Context) {

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

	// Retrieve role data from the database
	role, err := ReadByID(db.DB, domains.Roles{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving role data from the database. Error: ", err.Error())
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

	// Generate a user structure as a response
	details := RolesDetails{
		ID:          role.ID,
		Name:        role.Name,
		CompanyID:   role.OwningCompanyID,
		CompanyName: companyName,
		CreatedAt:   role.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateRole 		Handles the update of a role.
// @Summary        	Update role
// @Description    	Update one role.
// @Tags			Roles
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true		"Company ID"
// @Param			ID					path			string				true		"Role ID"
// @Param			request				body			roles.RolesIn		true		"Role query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses				"Invalid request"
// @Failure			401					{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses				"Forbidden"
// @Failure			500					{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/roles/{companyID}/{ID}	[put]
func (db Database) UpdateRole(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the role ID from the request parameter
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

	// Parse the incoming JSON request into a RoleIn struct
	role := new(RolesIn)
	if err := ctx.ShouldBindJSON(role); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the role with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Roles{}, objectID); err != nil {
		logrus.Error("Error checking if the role with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the role data in the database
	dbRole := &domains.Roles{
		Name: role.Name,
	}
	if err = domains.Update(db.DB, dbRole, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteRole	 	Handles the deletion of a role.
// @Summary        	Delete role
// @Description    	Delete one role.
// @Tags			Roles
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true		"Company ID"
// @Param			ID					path			string			true		"Role ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/roles/{companyID}/{ID}	[delete]
func (db Database) DeleteRole(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the role ID from the request parameter
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

	// Check if the role with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Roles{}, objectID); err != nil {
		logrus.Error("Error checking if the role with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the role data from the database
	if err := domains.Delete(db.DB, &domains.Roles{}, objectID); err != nil {
		logrus.Error("Error deleting role data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
