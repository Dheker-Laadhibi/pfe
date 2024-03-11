package interns

import (
	"labs/constants"
	"labs/domains"

	"net/http"
	"strconv"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

/**

IMPORTANT:
The user ID represents the unique identifier of an employee who holds the role of former intern groups.
Please ensure that appropriate permissions and access controls are in place for this user.
*/
//last update 26/02/2024 by dheker 21:27

// CreateIntern 	Handles the creation of a new intern.
// @Summary        	Create intern
// @Description    	Create a new intern.
// @Tags			Interns
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		   path			string				true	"companyID"
// @Param			request			body			interns.InternsIn	true	"Intern query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/interns/{companyID}	[post]
func (db Database) CreateIntern(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
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

	// Parse the incoming JSON request into a InternIn struct
	intern := new(InternsIn)
	if err := ctx.ShouldBindJSON(intern); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Hash the intern's password
	hash, _ := bcrypt.GenerateFromPassword([]byte(intern.Password), bcrypt.DefaultCost)

	// Create a new intern in the database
	dbIntern := &domains.Interns{
		ID:                         uuid.New(),
		Firstname:                  intern.Firstname,
		Lastname:                   intern.Lastname,
		Email:                      intern.Email,
		Password:                   string(hash),
		EducationalSupervisorEmail: intern.EducationalSupervisorEmail,
		EducationalSupervisorName:  intern.EducationalSupervisorName,
		EducationalSupervisorPhone: intern.EducationalSupervisorPhone,
		SupervisorID:               intern.SupervisorID,
		//CreatedByUserID: session.UserID,
		CompanyID: intern.CompanyID,
	}
	if err := domains.Create(db.DB, dbIntern); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

/*Gettttttttttttttt Alllllllllllllllllllllllllllllllll interns of one instructer*/

// ReadInterns 		Handles the retrieval of all interns.
// @Summary        	Get interns
// @Description    	Get all interns.
// @Tags			Interns
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		    path		string		true	"companyID"
// @Param			SupervisorID		    path		string		true	"SupervisorID"
// @Success			200				{object}	interns.InternsPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/interns/all/{companyID}/{SupervisorID}	[get]
func (db Database) ReadInterns(ctx *gin.Context) {

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

	// Parse and validate the user ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	SupervisorID, err := uuid.Parse(ctx.Param("SupervisorID"))
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

	// Check if the employee belongs to the specified user
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all intern data from the database
	interns, err := ReadAllPagination(db.DB, []domains.Interns{}, SupervisorID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all intern data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Interns{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a intern structure as a response
	response := InternsPagination{}
	dataTableIntern := []InternsTable{}
	for _, intern := range interns {

		dataTableIntern = append(dataTableIntern, InternsTable{
			ID:        intern.ID,
			Firstname: intern.Firstname,
			Lastname:  intern.Lastname,
			Email:     intern.Email,
		})
	}
	response.Items = dataTableIntern
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

/* Numbeeeeeeeeeeeeeeeeeer of Alll Internnnnnn in one company                                    */

// ReadInternsCount 	Handles the retrieval the number of all interns.
// @Summary        	Get number of  interns
// @Description    	Get number of all interns.
// @Tags			Interns
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true	"companyID "
// @Success			200						{object}		interns.InternsCount
// @Failure			400						{object}		utils.ApiResponses	"Invalid request"
// @Failure			401						{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses	"Forbidden"
// @Failure			500						{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/interns/count/{companyID}	[get]
// to add relation with company
func (db Database) ReadInternsCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
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
	interns, err := domains.ReadTotalCount(db.DB, &[]domains.Interns{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all intern data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a intern structure as a response
	internsCount := InternsCount{
		Count: interns,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, internsCount)
}

/*Geeeeeeeeeeeeet Onnnnnnnnnnnnnnnnnnnnnnnnnnnnne           inteeeeeeeeeeeeeeeeeeeeeeeeeeern*/

// ReadIntern 		Handles the retrieval of one intern.
// @Summary        	Get intern
// @Description    	Get one intern.
// @Tags			Interns
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			  true	   "company ID"
// @Param			ID					path			string			  true	    "internID "
// @Success			200					{object}		interns.InternsDetails
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/interns/{companyID}/{ID}	[get]
func (db Database) ReadIntern(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the intern ID from the request parameter
	ID, err := uuid.Parse(ctx.Param("ID"))
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

	// Retrieve intern data from the database
	intern, err := ReadByID(db.DB, domains.Interns{}, ID)
	if err != nil {
		logrus.Error("Error retrieving intern data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a intern structure as a response
	details := InternsDetails{
		ID:               intern.ID,
		Firstname:        intern.Firstname,
		Lastname:         intern.Lastname,
		Email:            intern.Email,
		LevelOfEducation: intern.LevelOfEducation,
		University:       intern.University,
		StartDate:        intern.StartDate,
	}
	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateIntern 		Handles the update of a intern.
// @Summary        	Update intern
// @Description    	Update intern.
// @Tags			Interns
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true	"companyID"
// @Param			ID					path			string				true	"Intern ID"
// @Param			request				body			interns.InternsIn		true	"Intern query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/interns/update/{companyID}/{ID}	[put]
func (db Database) UpdateIntern(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the intern ID from the request parameter
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

	// Parse the incoming JSON request into a InternIn struct
	intern := new(InternsIn)
	if err := ctx.ShouldBindJSON(intern); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the intern with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Interns{}, objectID); err != nil {
		logrus.Error("Error checking if the intern with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the intern data in the database
	dbIntern := &domains.Interns{
		Firstname: intern.Firstname,
		Lastname:  intern.Lastname,
		Email:     intern.Email,
	}
	if err = domains.Update(db.DB, dbIntern, objectID); err != nil {
		logrus.Error("Error updating intern data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteIntern	 	Handles the deletion of a intern.
// @Summary        	Delete intern
// @Description    	Delete one intern.
// @Tags			Interns
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"companyID"
// @Param			SupervisorID					path			string			true	"SupervisorID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/interns/Delete/{companyID}/{SupervisorID}	[delete]
func (db Database) DeleteIntern(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the intern ID from the request parameter
	SupervisorID, err := uuid.Parse(ctx.Param("SupervisorID"))
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

	// Check if the intern with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Interns{}, SupervisorID); err != nil {
		logrus.Error("Error checking if the intern with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the intern data from the database
	if err := domains.Delete(db.DB, &domains.Interns{}, SupervisorID); err != nil {
		logrus.Error("Error deleting intern data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())

}
