package mission_orders

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

// CreateMissionsOrders  Handles the creation of a new MissionOrders.
// @Summary        	Create MissionOrders
// @Description    	Create a new MissionOrders.
// @Tags			MissionOrders
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID	   path			string			true	"companyID"
// @Param			userID		   path			string			true	     "userID"
// @Param			request			body		mission_orders.MissionOrdersIn	true "MissionOrdersIn query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"MissionOrdersal Server Error"
// @Router			/missions/{companyID}/{userID}	[post]
func (db Database) CreateMissionOrders(ctx *gin.Context) {

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
	missionO := new(MissionOrdersIn)
	if err := ctx.ShouldBindJSON(missionO); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	layout := "2006-01-02" // Format de date année-mois-jour

	// Analyser la date ExpDate en tant que time.Time
	dt, err := time.Parse(layout, missionO.EndDate)
	if err != nil {
		logrus.Error("Erreur lors de l'analyse de la date : ", err.Error())
		// Gérer l'erreur ici
	}

	// Analyser la date ExpDate en tant que time.Time
	dtStart, err := time.Parse(layout, missionO.StartDate)
	if err != nil {
		logrus.Error("Erreur lors de l'analyse de la date : ", err.Error())
		// Gérer l'erreur ici
	}

	// Create a new MissionOrders  in the database
	dbMission := &domains.MissionOrders{
		ID:           uuid.New(),
		Object:       missionO.Object,
		Description:  missionO.Description,
		AdressClient: missionO.AdressClient,
		EndDate:      dt,
		StartDate:    dtStart,
		Transport:    missionO.Transport,
		UserID:       userID, //IDUser
		CompanyID: companyID,
	}

	if err := domains.Create(db.DB, dbMission); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadMissionOrders	Handles the retrieval of all MissionOrders.
// @Summary        		Get MissionsOrdes
// @Description    		Get all missions Orders.
// @Tags				MissionOrders
// @Produce				json
// @Security 			ApiKeyAuth
// @Param			    page			query		int			false		"Page"
// @Param			    limit			query		int			false		"Limit"
// @Param				companyID		path		string		true		"companyID"
// @Success				200					{array}			mission_orders.MissionOrdersDetails
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/missions/All/{companyID}	[get]
func (db Database) ReadMissionsOrders(ctx *gin.Context) {

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

	// Parse and validate the companyID from the request parameter
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

	// Check if the employee belongs to the specified user
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all user data from the database
	missions, err := ReadAllPagination(db.DB, []domains.MissionOrders{},session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	count, err := domains.ReadTotalCount(db.DB, &domains.MissionOrders{},"company_id", companyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a mission orders  structure as a response
	response := MissionsPagination{}
	dataTableMission := []MissionsTable{}
	for _, mission := range missions {

		dataTableMission = append(dataTableMission, MissionsTable{
			ID:          mission.ID,
			Object:      mission.Object,
			Description: mission.Description,
			StartDate:   mission.StartDate,
			Transport:   mission.Transport,
			EndDate:     mission.EndDate,
		})
	}
	response.Items = dataTableMission
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS,response)
}

// ReadMissionOrdersCount	Handles the retrieval the number of all MissionOrders.
// @Summary        			Get MissionOrders count
// @Description    			Get all mission_orders count.
// @Tags					MissionOrders
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					companyID				path			string		true		"companyID"
// @Success					200					{object}		mission_orders.MissionOrdersCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router					/missions/count/{companyID}	[get]
func (db Database) ReadMissionOrdersCount(ctx *gin.Context) {

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

	// Retrieve all MissionOrders data from the database
	mission_orders, err := domains.ReadTotalCount(db.DB, &domains.MissionOrders{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := MissionOrdersCount{
		Count: mission_orders,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}

// ReadPresence		    Handles the retrieval of one MissionOrders.
// @Summary        		Get MissionOrders
// @Description    		Get one MissionOrders.
// @Tags				MissionOrders
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID				path			string		true		"companyID"
// @Param				ID						path			string		true		"MissionOrders ID"
// @Success				200						{object}		mission_orders.MissionOrdersDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/missions/get/{companyID}/{ID}	[get]
func (db Database) ReadMissionOrders(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	
	// Parse and validate the companyID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the MissionOrders ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
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

	// Retrieve MissionOrders data from the database
	MissionOrders, err := ReadByID(db.DB, domains.MissionOrders{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving MissionOrders data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a MissionOrders  structure as a response
	details := MissionOrdersDetails{
		ID:        MissionOrders.ID,
		Object: MissionOrders.Object,
		Description: MissionOrders.Description,
		Transport: MissionOrders.Transport,
		EndDate:   MissionOrders.EndDate,
		StartDate: MissionOrders.StartDate,
		UserID:    MissionOrders.UserID,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdatePresence 	    Handles the update of a MissionOrders.
// @Summary        		Update MissionOrders
// @Description    		Update one MissionOrders.
// @Tags				MissionOrders
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID				path			string							true		"companyID "
// @Param				ID						path			string							true		"ID "
// @Param				request					body			mission_orders.MissionOrdersIn	true		"MissionOrdersIn query params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/missions/update/{companyID}/{ID}	[put]
func (db Database) UpdateMissionOrders(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID  from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the MissionOrders ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
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

	// Parse the incoming JSON request into a MissionOrdersIn struct
	MissionOrders := new(MissionOrdersIn)
	if err := ctx.ShouldBindJSON(MissionOrders); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the MissionOrders with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.MissionOrders{}, objectID); err != nil {
		logrus.Error("Error checking if the MissionOrders with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the MissionOrders data in the database
	dbMissionOrders := &domains.MissionOrders{
		Description: MissionOrders.Description,
		Object:      MissionOrders.Object,
		Transport: MissionOrders.Transport,
	}
	if err = domains.Update(db.DB, dbMissionOrders, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteMissionOrders 	Handles the deletion of a MissionOrders.
// @Summary        		Delete MissionOrders
// @Description    		Delete one MissionOrders.
// @Tags				MissionOrders
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID				path			string			true		"companyID"
// @Param				ID						path			string			true		"MissionOrders ID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"MissionOrderseal Server Error"
// @Router				/missions/delete/{companyID}/{ID}	[delete]
func (db Database) 	DeleteMissionOrders(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the MissionOrders ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
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

	// Delete the MissionOrders data from the database
	if err := domains.Delete(db.DB, &domains.MissionOrders{}, objectID); err != nil {
		logrus.Error("Error deleting user data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
