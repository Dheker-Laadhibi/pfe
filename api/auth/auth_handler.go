package auth

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// SignupUser 		Handles the user signup process.
// @Summary        	Signup
// @Description    	Create a new user account.
// @Tags			Authentification
// @Accept			json
// @Produce			json
// @Param			request		body		auth.Signup		true	"User query params"
// @Success			201			{object}	utils.ApiResponses
// @Failure			400			{object}	utils.ApiResponses		"Invalid request"
// @Failure			401			{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403			{object}	utils.ApiResponses		"Forbidden"
// @Failure			500			{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/auth/signup	[post]
func (db Database) SignupUser(ctx *gin.Context) {

	// Parse the incoming JSON request into a Signup struct
	user := new(Signup)
	if err := ctx.ShouldBindJSON(user); err != nil {
		logrus.Error("Failed to map request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check for email duplication
	if check := CheckEmailDuplication(db.DB, user.Email); check != nil {
		logrus.Error("Error checking email duplication. Error: ", check)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new company in the database
	dbCompany := &domains.Companies{
		ID:   uuid.New(),
		Name: user.CompanyName,
	}
	if err := domains.Create(db.DB, dbCompany); err != nil {
		logrus.Error("Error creating a new company in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Hash the user's password
	// becrypt + bcrypt.DefaultCost salt random caracter  in pwd 
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Create a new user in the database
	dbUser := &domains.Users{
		ID:        uuid.New(),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  string(hash),
		Status:    true,
		CompanyID: dbCompany.ID,
	}
	if err := domains.Create(db.DB, dbUser); err != nil {
		logrus.Error("Error creating a new user in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Create default role for the user
	dbRole := &domains.Roles{
		ID:              uuid.New(),
		Name:            constants.DEFAULT_ROLE,
		OwningCompanyID: dbCompany.ID,
		CreatedByUserID: dbUser.ID,
	}
	if err := domains.Create(db.DB, dbRole); err != nil {
		logrus.Error("Error creating a default role for the user. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Assign the role to the user
	dbUserRole := &domains.UsersRoles{
		UserID:    dbUser.ID,
		RoleID:    dbRole.ID,
		CompanyID: dbCompany.ID,
	}
	if err := domains.Create(db.DB, dbUserRole); err != nil {
		logrus.Error("Error assigning the role to the user. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// SigninUser 		Handles the user signin process.
// @Summary			Signin
// @Description		Authenticate and log in a user.
// @Tags			Authentification
// @Accept			json
// @Produce			json
// @Param			request		body		auth.Signin		true	"User query params"
// @Success			200			{object}	auth.LoggedInResponse
// @Failure			400			{object}	utils.ApiResponses		"Invalid request"
// @Failure			401			{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403			{object}	utils.ApiResponses		"Forbidden"
// @Failure			500			{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/auth/signin	[post]
func (db Database) SigninUser(ctx *gin.Context) {

	// Parse the incoming JSON request into a Signin struct
	user := new(Signin)
	if err := ctx.ShouldBindJSON(user); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve user data by email
	data, err := ReadByEmailActive(db.DB, user.Email)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Compare the entered password with the stored password
	if isTrue := utils.ComparePassword(data.Password, user.Password); !isTrue {
		logrus.Error("Password comparison failed.")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNAUTHORIZED, utils.Null())
		return
	}

	// Retrieve user roles from the database
	//dbRoles tous les roles de db 	
	dbRoles, err := domains.ReadUsersRoles(db.DB, data.ID, data.CompanyID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Prepare roles for the JWT session
	 roles := []domains.RolesSessionJWT{}
	 // parcourt tous les rôles de l'utilisateur qui ont été récupérés de la base de données.
	  for _, role := range dbRoles {
       //obtenir le nom du rôle associé à cet identifiant.
		name, err := domains.ReadRoleName(db.DB, role.RoleID)
		if err != nil {
			logrus.Error("Error retrieving data from the database. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
			return
		}
// nom role  ajouté à la liste des rôles de session JWT.
		roles = append(roles, domains.RolesSessionJWT{
			ID:        role.RoleID,
			Name:      name,
			CompanyID: role.CompanyID,
		})
	}

	// Generate JWT token
	token := utils.GenerateToken(data.ID, data.CompanyID, roles)
// une réponse structurée contenant à la fois le token JWT + info user 
	// Prepare the response
	response := LoggedInResponse{
		AccessToken: token,
		User: LoggedIn{
			ID:             data.ID,
			Name:           data.Firstname + " " + data.Lastname,
			Email:          data.Email,
			ProfilePicture: data.ProfilePicture,
			CompanyID:      data.CompanyID,
		},
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}
