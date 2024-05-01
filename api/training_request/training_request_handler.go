package training_request

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

//create TrainingRequest

// CreateTrainingRequest  Handles the creation of a new TrainingRequest.
// @Summary        	Create TrainingRequest
// @Description    	Create a new TrainingRequest.
// @Tags			TrainingRequest
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID	    path			string				true	"companyID"
// @Param			userID		    path			string				true	"userID"
// @Param			request			body		    training_request.TrainingRequestIn	true "TrainingRequestIn query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Training request  Server Error"
// @Router			/training_request/{companyID}/{userID}	[post]
func (db Database) CreateTrainingRequestByUser(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user id  from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
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

	// indicating if the exact user making the request is the same as the user associated with the session.
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a PresenceIn struct
	trainingR := new(TrainingRequestIn)
	if err := ctx.ShouldBindJSON(trainingR); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	layout := "2006-01-02" // Format de date année-mois-jour

	// Analyser la date ExpDate en tant que time.Time
	dt, err := time.Parse(layout, trainingR.RequestDate)
	if err != nil {
		logrus.Error("Erreur lors de l'analyse de la date : ", err.Error())
		// Gérer l'erreur ici
	}

	// Create a new training request  in the database
	dbTraining := &domains.TrainingRequest{
		ID:            uuid.New(),
		TrainingTitle: trainingR.TrainingTitle,
		Description:   trainingR.Description,
		Reason:        trainingR.Reason,
		RequestDate:   dt,
		UserID:        userID, //IDUser
		CompanyID: companyID,
	}

	if err := domains.Create(db.DB, dbTraining); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())

}

// ReadTrainingRequests	Handles the retrieval of all TrainingRequests.
// @Summary        		Get TrainingRequests
// @Description    		Get all  training request .
// @Tags				TrainingRequest
// @Produce				json
// @Security 			ApiKeyAuth
// @Param			    page			    query		    int			false		"Page"
// @Param			    limit			    query		    int			false		"Limit"
// @Param				companyID		    path			string		true		"companyID "
// @Success				200					{array}			training_request.TrainingRequestPagination
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure				500					{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/training_request/All/{companyID}	[get]
func (db Database) ReadTrainingsRequest(ctx *gin.Context) {

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



	// Parse and validate the  company ID from the request parameter
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

	// Retrieve all user data from the database
	Trainings, err := ReadAllPagination(db.DB, []domains.TrainingRequest{}, companyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all TrainingRequest data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	count, err := domains.ReadTotalCount(db.DB, &domains.TrainingRequest{}, "company_id", companyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	

	// Generate a  training request  orders  structure as a response
	response := TrainingRequestPagination{}
	dataTableTraining := []TrainingRequestTable{}
	for _, Training := range Trainings {

		dataTableTraining = append(dataTableTraining, TrainingRequestTable{
			ID:            Training.ID,
			TrainingTitle: Training.TrainingTitle,
			Description:   Training.Description,
			RequestDate:   Training.RequestDate,
			Reason:        Training.Reason,
			UserID:        Training.UserID,
		})
	}
	response.Items = dataTableTraining
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}


// ReadTrainingsRequestCount	Handles the retrieval the number of all TrainingRequest.
// @Summary        			Get TrainingRequest count
// @Description    			Get all TrainingRequest count.
// @Tags					TrainingRequest
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					companyID			path			string		true		"companyID "
// @Success					200					{object}		training_request.TrainingRequestsCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router					/training_request/count/{companyID}	[get]
func (db Database) ReadTrainingsCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the TrainingRequests ID from the request parameter
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


	// Retrieve all  training request  data from the database
	training_request, err := domains.ReadTotalCount(db.DB, &domains.TrainingRequest{}, "company_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := TrainingRequestsCount{
		Count: training_request,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}



// ReadPresence			Handles the retrieval of one TrainingRequests.
// @Summary        		Get TrainingRequests
// @Description    		Get one TrainingRequests.
// @Tags				TrainingRequest
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID				path			string		true		"companyID"
// @Param				ID						path			string		true		"training_request ID"
// @Success				200						{object}		training_request.TrainingRequestDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router				/training_request/get/{companyID}/{ID}	[get]
func (db Database) ReadTrainingsRequests(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the training_request ID from the request parameter
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

	// Retrieve training_request data from the database
	trainings, err := ReadByID(db.DB, domains.TrainingRequest{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving training_request data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a training request   structure as a response
	details := TrainingRequestDetails{
		ID:              trainings.ID,
		TrainingTitle:   trainings.TrainingTitle,
		Reason:          trainings.Reason,
		Description:     trainings.Description,
		RequestDate:     trainings.RequestDate,
		UserID:          trainings.UserID,
		DecisionCompany: trainings.DecisionCompany,
	
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdatePresence 	    Handles the update of a TrainingsRequests.
// @Summary        		Update TrainingRequest
// @Description    		Update one TrainingRequest.
// @Tags				TrainingRequest
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID				path		    string		                    true	    "companyID"
// @Param				ID						path			string							true	    "TrainingRequestIn ID"
// @Param				request					body			training_request.TrainingRequestDescision	true		"TrainingRequestIn query params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Presenceal Server Error"
// @Router			    /training_request/update/{companyID}/{ID}	[put]
func (db Database) UpdateTraining(ctx *gin.Context) {
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

	// Parse and validate the training_request ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a training_request struct
	TraininIn := new(TrainingRequestDescision)
	if err := ctx.ShouldBindJSON(TraininIn); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the training_request with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.TrainingRequest{}, objectID); err != nil {
		logrus.Error("Error checking if the TraininIn  with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the training_request data in the database
	dbTraining := &domains.TrainingRequest{
		DecisionCompany: TraininIn.DecisionCompany,
	}

	if err = domains.Update(db.DB, dbTraining, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
	
}

// DeleteTrainingRequest 	Handles the deletion of a TrainingRequest.
// @Summary        		Delete TrainingRequest
// @Description    		Delete one TrainingRequest.
// @Tags				TrainingRequest
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				companyID				path		    string		    true	    "companyID"
// @Param				ID						path			string			true		"TrainingRequest ID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"MissionOrderseal Server Error"
// @Router				/training_request/delete/{companyID}/{ID}	[delete]
func (db Database) DeleteTrainingsRequest(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the companyID  from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the training_request ID from the request parameter
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

	// Delete the training_request data from the database
	if err := domains.Delete(db.DB, &domains.TrainingRequest{}, objectID); err != nil {
		logrus.Error("Error deleting training request  data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
