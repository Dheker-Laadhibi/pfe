package auth

import "github.com/google/uuid"

// @Description	Signup represents the account information for signing up.
type Signup struct {
	Firstname   string `json:"firstName" binding:"required,min=3,max=30"`    // Firstname is the first name of the user. It should be between 3 and 30 characters.
	Lastname    string `json:"lastName" binding:"required,min=3,max=35"`     // Lastname is the last name of the user. It should be between 3 and 35 characters.
	Email       string `json:"email" binding:"required,email,max=255"`       // Email is the email address of the user. It is required, should be a valid email, and maximum length is 255 characters.
	Password    string `json:"password" binding:"required,min=10,max=255"`   // Password is the user's password. It is required, and its length should be between 10 and 255 characters.
	CompanyName string `json:"companyName" binding:"required,min=2,max=255"` // CompanyName is the name of the user's company. It is required and should be between 2 and 255 characters.
} //@name Signup

// @Description	Signin represents the information required for signing in.
type Signin struct {
	Email    string `json:"email" binding:"required,email,max=255"`     // Email is the email address of the user. It is required, should be a valid email, and maximum length is 255 characters.
	Password string `json:"password" binding:"required,min=10,max=255"` // Password is the user's password. It is required, and its length should be between 10 and 255 characters.
} //@name Signin

// @Description	LoggedInResponse represents the response structure after successful login.
type LoggedInResponse struct {
	AccessToken string   `json:"accessToken"` // AccessToken is the token obtained after successful login for authentication purposes.
	User        LoggedIn `json:"user"`        // User is the structure containing details of the logged-in user.
} //@name LoggedInResponse

// @Description	LoggedIn represents the user details after successful login.
type LoggedIn struct {
	ID             uuid.UUID `json:"ID"`             // ID is the unique identifier for the user.
	Name           string    `json:"name"`           // Name is the name of the user.
	Email          string    `json:"email"`          // Email is the email address of the user.
	ProfilePicture string    `json:"profilePicture"` // ProfilePicture is the URL or path to the user's profile picture.
	CompanyID      uuid.UUID `json:"workCompanyId"`  // CompanyID is the unique identifier for the user's company.
} //@name LoggedIn
