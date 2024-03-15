package projects

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateProject		Handles the creation of a new Project.
// @Summary        	Create Project
// @Description    	Create a new Project.
// @Tags			Project
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				true		"Company ID"
// @Param			request			body			projects.ProjectIn		true		"project query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/projects/{companyID}	[post]
func (db Database) CreateProject(ctx *gin.Context) {

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

	// Parse the incoming JSON request into a ProjectIn struct
	project := new(ProjectIn)
	if err := ctx.ShouldBindJSON(project); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	layout := "2006-01-02" // Format de date année-mois-jour

	// Analyser la date ExpDate en tant que time.Time
	dt, err := time.Parse(layout, project.ExpDate)
	if err != nil {
		logrus.Error("Erreur lors de l'analyse de la date : ", err.Error())
		// Gérer l'erreur ici
	}

	// Create a new project in the database
	dbProject := &domains.Project{
		ID:           uuid.New(),
		Code:         project.Code,
		Projectname:  project.Projectname,
		Technologies: project.Technologies,
		Description:  project.Description,
		ExpDate:      dt,
		CompanyID:    project.CompanyID,
		// ExpDate: project.,
	}

	// Vérifier si le code de projet existe déjà dans la base de données
	exists, err := domains.CheckProjectCodeExists(db.DB, project.Code)
	if err != nil {
		logrus.Error("Erreur lors de la vérification de l'existence du code de projet dans la base de données : ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	} else if exists {
		logrus.Error("Le code de projet existe déjà dans la base de données.")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DUPLICATE_CODE, utils.Null())
		return
	}

	// Créer le projet uniquement si le code de projet n'existe pas déjà dans la base de données
	if err := domains.Create(db.DB, dbProject); err != nil {
		logrus.Error("Erreur lors de l'enregistrement des données dans la base de données. Erreur: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Répondre avec succès
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())

}

// ReadProjects		Handles the retrieval of all Project.
// @Summary        	Get projects
// @Description    	Get all projects.
// @Tags			Project
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		path		string		true		"Company ID"
// @Success			200				{object}	projects.projectsPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/projects/{companyID}	[get]
func (db Database) ReadProjects(ctx *gin.Context) {

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

	// Retrieve all project data from the database
	projects, err := ReadAllPagination(db.DB, []domains.Project{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all project data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Project{}, "company_id", companyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a project structure as a response
	response := projectsPagination{}
	dataTableProject := []ProjectTable{}
	for _, project := range projects {

		dataTableProject = append(dataTableProject, ProjectTable{
			ID:           project.ID,
			Code:         project.Code,
			Projectname:  project.Projectname,
			CompanyID:    project.CompanyID,
			Technologies: project.Technologies,
			ExpDate:      project.ExpDate,
		})
	}
	response.Items = dataTableProject
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadProjectsList 	Handles the retrieval the list of all Projects.
// @Summary        	Get list of  Projects
// @Description    	Get list of all Projects.
// @Tags			Project
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Success			200					{array}			projects.ProjectsList
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/projects/{companyID}/list	[get]
func (db Database) ReadProjectsList(ctx *gin.Context) {

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

	// Retrieve all project data from the database
	projects, err := ReadAllList(db.DB, []domains.Project{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all project data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a project structure as a response
	projectList := []ProjectsList{}
	for _, project := range projects {
		projectList = append(projectList, ProjectsList{
			ID:          project.ID,
			Code:        project.Code,
			Projectname: project.Projectname,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, projectList)
}

// ReadProjectsCount 	Handles the retrieval the number of all projects.
// @Summary        	Get number of  projects
// @Description    	Get number of all projects.
// @Tags			Project
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true	"Company ID"
// @Success			200						{object}		projects.ProjectsCount
// @Failure			400						{object}		utils.ApiResponses	"Invalid request"
// @Failure			401						{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses	"Forbidden"
// @Failure			500						{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/projects/{companyID}/count	[get]
func (db Database) ReadProjectsCount(ctx *gin.Context) {

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

	// Retrieve all project data from the database
	projects, err := domains.ReadTotalCount(db.DB, &[]domains.Project{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all project data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a project structure as a response
	ProjectsCount := ProjectsCount{
		Count: projects,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, ProjectsCount)
}

// ReadProject		Handles the retrieval of one Project.
// @Summary        	Get Project
// @Description    	Get one Project.
// @Tags			Project
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"Project ID"
// @Success			200					{object}		projects.projectsDetails
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/projects/{companyID}/{ID}	[get]
func (db Database) ReadProject(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the project ID from the request parameter
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

	// Retrieve project data from the database
	project, err := ReadByID(db.DB, domains.Project{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving project data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a project structure as a response
	details := projectsDetails{
		ID:           project.ID,
		Code:         project.Code,
		Projectname:  project.Projectname,
		CompanyID:    project.CompanyID,
		Technologies: project.Technologies,
		ExpDate:      project.ExpDate,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateProject 		Handles the update of a Project .
// @Summary        	Update Project
// @Description    	Update Project .
// @Tags			Project
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true	"Company ID"
// @Param			ID					path			string				true	"Project  ID"
// @Param			request				body			projects.ProjectIn		true	"Project query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/projects/{companyID}/{ID}	[put]
func (db Database) UpdateProject(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the project ID from the request parameter
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

	// Parse the incoming JSON request into a ProjectIn struct
	project := new(ProjectIn)
	if err := ctx.ShouldBindJSON(project); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the project with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Project{}, objectID); err != nil {
		logrus.Error("Error checking if the project with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the project data in the database
	dbProject := &domains.Project{

		Projectname:  project.Projectname,
		Technologies: project.Technologies,
	}
	if err = domains.Update(db.DB, dbProject, objectID); err != nil {
		logrus.Error("Error updating project data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteProject	 	Handles the deletion of a Project.
// @Summary        	Delete Project
// @Description    	Delete one Project.
// @Tags			Project
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"Project ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/projects/{companyID}/{ID}	[delete]
func (db Database) DeleteProject(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the project ID from the request parameter
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

	// Check if the project with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Project{}, objectID); err != nil {
		logrus.Error("Error checking if the project with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the project data from the database
	if err := domains.Delete(db.DB, &domains.Project{}, objectID); err != nil {
		logrus.Error("Error deleting project data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// AssignProject	AssignProjectToCondidats.
// @Summary        	assign Project
// @Description    	assign one Project.
// @Tags			Project
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"candidate ID"
// @Param			request			     body			projects.codeProject		true  "project query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/projects/{companyID}/{ID}/assign	[post]
func (db Database) AssignProjectToCondidats(ctx *gin.Context) {

	// Extract JWT values from the context
	//session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Check if the candidate with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Companies{}, companyID); err != nil {
		logrus.Error("Error checking if the candidate with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Parse and validate the candidate ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the candidate with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Condidats{}, objectID); err != nil {
		logrus.Error("Error checking if the candidate with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	/*	// Parse the incoming JSON request into a ProjectIn struct
		project := new(ProjectIn)
		if err := ctx.ShouldBindJSON(project); err != nil {
			logrus.Error("Error mapping request from frontend. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
			return
		}
	*/
	projectCode := new(codeProject)
	if err := ctx.ShouldBindJSON(projectCode); err != nil {
		logrus.Error("Error parsing request body. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Error parsing request body.", nil)
		return
	}

	// Query database to find project ID based on project code
	projectID, err := domains.FindProjectIDByCode(db.DB, projectCode.Code)
	if err != nil {
		logrus.Error("Error finding project ID by project code. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Error finding project ID by project code.", nil)
		return
	}
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		logrus.Error("Error parsing project ID: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Error parsing project ID.", nil)
		return
	}

	// Assign candidate to project
	if err := domains.AssignProjectCondidat(db.DB, parsedProjectID, objectID, companyID); err != nil {
		logrus.Error("Error assigning candidate to project. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Error assigning candidate to project.", nil)
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, "Project assigned to candidate successfully.", nil)
}
