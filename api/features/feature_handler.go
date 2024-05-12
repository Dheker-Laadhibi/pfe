package features

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




// CreateFeature 	Handles the creation of a new feature.
// @Summary        	Create feature
// @Description    	Create a new feature.
// @Tags			Feature
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				true		"Company ID"
// @Param			request			body			features.FeatureIn		true	"Feature query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/features/{companyID}	[post]
func (db Database) CreateFeature(ctx *gin.Context) {

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
		logrus.Error("Error verifying employee belonging to company . Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a featureIn struct
	feature := new(FeatureIn)
	if err := ctx.ShouldBindJSON(feature); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new feature in the database
	dbfeature := &domains.Feature{
		ID:              uuid.New(),
		Featurename:            feature.Featurename,
          CompanyID: companyID,
		
	}
	if err := domains.Create(db.DB, dbfeature); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}



// ReadFeature 		Handles the retrieval of all features.
// @Summary        	Get features
// @Description    	Get all features.
// @Tags			Feature
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page				query		int				false		"Page"
// @Param			limit				query		int				false		"Limit"
// @Param			companyID			path		string			true		"Company ID"
// @Success			200					{object}	features.FeaturePagination
// @Failure			400					{object}	utils.ApiResponses			"Invalid request"
// @Failure			401					{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}	utils.ApiResponses			"Forbidden"
// @Failure			500					{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/features/{companyID}	[get]
func (db Database) ReadFeatures(ctx *gin.Context) {

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

	// Retrieve all feature data from the database
	features, err := ReadAllPagination(db.DB, []domains.Feature{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Feature{}, "company_id", session.CompanyID)
	if err != nil {

		logrus.Error("Error occurred while finding total  count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Generate a user structure as a response
	response := FeaturePagination{}
	dataTableFeature := []FeatureTable{}
	for _, feature := range features {

		dataTableFeature = append(dataTableFeature,FeatureTable{
			ID:        feature.ID,
		Featurename:      feature.Featurename,
		
		})
	}
	response.Items = dataTableFeature
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}





// ReadFeaturesList	Handles the retrieval the list of all features.
// @Summary        	Get list of  features
// @Description    	Get list of all features.
// @Tags			Feature
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			   	path			string			true	"Company ID"
// @Success			200						{array}			features.FeatureTable
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/features/{companyID}/list	[get]
func (db Database) ReadFeaturesList(ctx *gin.Context) {

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

	// Retrieve all feature data from the database
	features, err := ReadAllList(db.DB, []domains.Feature{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a feature structure as a response
	featureList := []FeatureTable{}
	for _, feature := range features {

		featureList = append(featureList, FeatureTable{
			ID:   feature.ID,
			Featurename: feature.Featurename,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, featureList)
}

// ReadFeaturesCount	Handles the retrieval the number of all features.
// @Summary        	Get number of  features
// @Description    	Get number of all features.
// @Tags			Feature
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true		"Company ID"
// @Success			200						{object}		features.FeatureCount
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/features/{companyID}/count	[get]
func (db Database) ReadFeaturesCount(ctx *gin.Context) {

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

	// Retrieve all feature data from the database
	features, err := domains.ReadTotalCount(db.DB, &[]domains.Feature{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a feature structure as a response
	featuresCount := FeatureCount{
		Count: features,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, featuresCount)
}

// Readfeature 		Handles the retrieval of one feature.
// @Summary        	Get feature
// @Description    	Get one feature.
// @Tags			Feature
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string			true	"Company ID"
// @Param			ID						path			string			true	"feature ID"
// @Success			200						{object}		features.FeatureTable
// @Failure			400						{object}		utils.ApiResponses		"Invalid request"
// @Failure			401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses		"Forbidden"
// @Failure			500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/features/{companyID}/{ID}	[get]
func (db Database) ReadFeature(ctx *gin.Context) {

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

	// Retrieve feature data from the database
	feature, err := ReadByID(db.DB, domains.Feature{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving feature data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	

	// Generate a user structure as a response
	details := FeatureTable{
		ID:          feature.ID,
		Featurename: feature.Featurename,
		
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateFeature	Handles the update of a feature.
// @Summary        	Update feature
// @Description    	Update one feature.
// @Tags			Feature
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true		"Company ID"
// @Param			ID					path			string				true		"feature ID"
// @Param			request				body			features.FeatureIn		true	"feature query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses				"Invalid request"
// @Failure			401					{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses				"Forbidden"
// @Failure			500					{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/features/{companyID}/{ID}	[put]
func (db Database) UpdateFeature(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the feature ID from the request parameter
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

	// Parse the incoming JSON request into featureIn struct
	feature := new(FeatureIn)
	if err := ctx.ShouldBindJSON(feature); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the feature with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Feature{}, objectID); err != nil {
		logrus.Error("Error checking if the feature with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the feature data in the database
	dbfeature := &domains.Feature{
		Featurename: feature.Featurename,
	}
	if err = domains.Update(db.DB, dbfeature, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteFeature 	Handles the deletion of a feature.
// @Summary        	Delete feature
// @Description    	Delete one feature.
// @Tags			Feature
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true		"Company ID"
// @Param			ID					path			string			true		"feature ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/features/{companyID}/{ID}	[delete]
func (db Database) DeleteFeature(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the feature ID from the request parameter
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

	// Check if the feature with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Feature{}, objectID); err != nil {
		logrus.Error("Error checking if the feature with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the feature data from the database
	if err := domains.Delete(db.DB, &domains.Feature{}, objectID); err != nil {
		logrus.Error("Error deleting feature data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
