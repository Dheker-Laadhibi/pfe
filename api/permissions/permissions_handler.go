package permissions

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



// CreateMissionsOrders  Handles the creation of a new permission.
// @Summary        	Create permission
// @Description    	Create a new permission.
// @Tags			Permissions
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID	   path			string			true	"companyID"
// @Param			roleID		   path			string			true	 "roleID"
// @Param			featureID	   path			string			true	"featureID"
// @Param			request			body		permissions.PermissionIn	true "PermissionIn query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"MissionOrdersal Server Error"
// @Router			/permissions/{companyID}/{featureID}/{roleID}	[post]
func (db Database) CreatePermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	
		// Parse and validate the user ID from the request parameter
		roleID, err := uuid.Parse(ctx.Param("roleID"))
		if err != nil {
			logrus.Error("Error  roleID mapping request from frontend . Invalid UUID format. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
			return
		}
		// Parse and validate the user ID from the request parameter
		featureID, err := uuid.Parse(ctx.Param("featureID"))
		if err != nil {
			logrus.Error("Error mapping request from frontend featureID. Invalid UUID format. Error: ", err.Error())
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
	permission := new(PermissionIn)
	if err := ctx.ShouldBindJSON(permission); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	

	// Create a new MissionOrders  in the database
	dbpermission := &domains.Permissions{
		ID:           uuid.New(),
	FeatureName: permission.FeatureName,
	ReadPerm: permission.ReadPerm,
	CreatePerm: permission.CreatePerm,
	UpdatePerm: permission.UpdatePerm,
	DeletePerm: permission.DeletePerm,
		FeatureID:    featureID,
		RoleID:       roleID, 
		CompanyID: companyID,
	}

	if err := domains.Create(db.DB, dbpermission); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}



// ReadPermissions	Handles the retrieval of all permissions.
// @Summary        	Get permissions
// @Description    	Get all permissions.
// @Tags			Permissions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false	"Page"
// @Param			limit			query		int			false	"Limit"
// @Param			companyID		path		string		true	"companyID"
// @Success			200				{object}	permissions.PermissionsPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/permissions/{companyID}	[get]
func (db Database) ReadPermissions(ctx *gin.Context) {

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

	// Parse and validate the companyID  from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the intern's value is among the allowed choices
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

	// Check if the employee belongs to the specified intenr
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all intern data from the database
	permissions, err := ReadAllPagination(db.DB, []domains.Permissions{}, companyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all intern data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count by one supervisor
	count, err := domains.ReadTotalCount(db.DB, &domains.Permissions{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a intern structure as a response
	response := PermissionsPagination{}
	dataTablePermission := []PermissionTable{}
	for _, permission := range permissions {

		dataTablePermission = append(dataTablePermission, PermissionTable{
			ID:          permission.ID,
			FeatureName: permission.FeatureName,
			CreatePerm:  permission.CreatePerm,
			UpdatePerm:  permission.UpdatePerm,
			ReadPerm:    permission.ReadPerm,
			DeletePerm:  permission.DeletePerm,
		})
	}
	response.Items = dataTablePermission
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadPermissionsCount Handles the retrieval the number of all permissions.
// @Summary        	Get number of  permissions
// @Description    	Get number of all permissions.
// @Tags			Permissions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true	"companyID "
// @Success			200						{object}		permissions.permissionCount
// @Failure			400						{object}		utils.ApiResponses	"Invalid request"
// @Failure			401						{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses	"Forbidden"
// @Failure			500						{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/permissions/{companyID}/count	[get]
// to add relation with company
func (db Database) ReadPermissionCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID  from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Retrieve all intern data from the database
	permissions, err := domains.ReadTotalCount(db.DB, &[]domains.Permissions{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all permissions data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a intern structure as a response
	permissionCount := permissionCount{
		Count: permissions,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, permissionCount)
}

// ReadPermission	Handles the retrieval of one permission.
// @Summary        	Get permission
// @Description    	Get one permission.
// @Tags			Permissions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path	string			  true	   "company ID"
// @Param			permissionID		path	string			  true	    "permissionID "
// @Success			200					{object}		permissions.PermissionsDetails
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/permissions/{companyID}/permission/{permissionID}	[get]
func (db Database) ReadPermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID  from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the intern ID from the request parameter
	permissionID, err := uuid.Parse(ctx.Param("permissionID"))
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

	// Retrieve permission data from the database
	permission, err := ReadByID(db.DB, domains.Permissions{}, permissionID)
	if err != nil {
		logrus.Error("Error retrieving intern data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a intern structure as a response
	details := PermissionsDetails{
		ID:         permission.ID,
		FeatureName: permission.FeatureName,
		CreatePerm: permission.CreatePerm,
		ReadPerm:   permission.ReadPerm,
		DeletePerm: permission.DeletePerm,
		UpdatePerm: permission.UpdatePerm,			
	}
	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdatePermission	Handles the update of a permission.
// @Summary        	Update permission
// @Description    	Update permission.
// @Tags			Permissions
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true	"companyID"
// @Param			permissionID	    path		    string				true	"permissionID"
// @Param			request				body			permissions.PermissionIn	true	"permission query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/permissions/{companyID}/update/{permissionID}	[put]
func (db Database) UpdatePermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID  from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the intern ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("permissionID"))
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

	// Parse the incoming JSON request into a InternIn struct
	permission := new(PermissionIn)
	if err := ctx.ShouldBindJSON(permission); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the intern with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Permissions{}, objectID); err != nil {
		logrus.Error("Error checking if the intern with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the intern data in the database
	dbpermission := &domains.Permissions{
		FeatureName: permission.FeatureName,
		CreatePerm:  permission.CreatePerm,
		UpdatePerm:  permission.UpdatePerm,
		ReadPerm: permission.ReadPerm,
		DeletePerm:  permission.DeletePerm,
	}
	if err = domains.Update(db.DB, dbpermission, objectID); err != nil {
		logrus.Error("Error updating intern data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
// DeletePermission  Handles the deletion of a permission.
// @Summary        	Delete permission
// @Description    	Delete one permission.
// @Tags			Permissions
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"companyID"
// @Param			permissionID		 path			string			true	"permissionID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/permissions/{companyID}/delete/{permissionID}	 [delete]
func (db Database) DeletePermission(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company  ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the intern ID from the request parameter
	permissionID, err := uuid.Parse(ctx.Param("permissionID"))
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

	

	// Delete the intern data from the database
	if err := domains.Delete(db.DB, &domains.Permissions{}, permissionID); err != nil {
		logrus.Error("Error deleting intern data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())







}
