package candidats

import (
	"github.com/google/uuid"
)

// @Description	CondidatIn represents the input structure for creating a new condidat.
type CondidatIn struct {
	Firstname string `json:"firstName" binding:"min=3,max=30"`	// Firstname is the first name of the user. It is required and should be between 3 and 30 characters.
	Lastname       string    `json:"lastName"  binding:"min=3,max=30"`  // Lastname is the last name of the user. It is required and should be between 3 and 35 characters.
	Adress         string    `json:"adress"  binding:"min=3,max=30"`
	Email          string    `json:"email" binding:"required,email,max=255"` // Email is the email address of the user. It is required, should be a valid email, and maximum length is 255 characters.
	Educationlevel string    `json:"education_level"  binding:"min=3,max=100"`
	RoleName        string    `json:"role_name"`   //  define condidat role 
	University     string    `json:"university"  binding:"min=3,max=30"`
	Password       string    `json:"password" binding:"required,min=10,max=255"` // Password is the user's password. It is required, and its length should be between 10 and 255 characters.

} //@name CondidatIn

// @Description	CondidtasPagination represents the paginated list of Condidats.
type CondidtasPagination struct {
	Items      []CondidatsTable `json:"items"`      // Items is a slice containing individual condidat details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of condidats in the entire list.
} //@name CondidtasPagination

// @Description	CondidatsTable represents a single condidat entry in a table.
type CondidatsTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the  condidat.
	Firstname string    `json:"firstname"` // Firstname is the first name of the condidat.
	Lastname  string    `json:"lastname"`  // Lastname is the last name of the condidat.
	Email     string    `json:"email"` 
	Adress         string  `json:"adress"`
	Educationlevel string `json:"educationlevel"`
	University     string`json:"university"`  
	Status   bool       `json:"status"`
} //@name CondidatsTable

// @Description	CondidatsList represents a simplified version of the Condidats for listing purposes.
type CondidatsList struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the condidat.
	Firstname string    `json:"firstname"` // Name is the name of the condidat.
	Lastname  string    `json:"lastname"`  // Name is the name of the condidat.
	Adress         string  `json:"adress"`
	Educationlevel string `json:"educationlevel"`
	University     string `json:"university"`
	Status   bool     `json:"status"`
} //@name CondidatsList

// @Description	CondidatsCount represents the count of condidats.
type CondidatsCount struct {
	Count uint `json:"count"` // Count is the number of condidats.
} //@name CondidatsCount

// @Description	CondidatDetails represents detailed information about a specific condidat.
type CondidatDetails struct {
	ID               uuid.UUID `json:"id"`          // ID is the unique identifier for the condidat.
	Firstname        string    `json:"firstname"`   // Name is the name of the condidat.
	Lastname         string    `json:"lastname"`    // last name of condidat
	CompanyID        uuid.UUID `json:"companyID"`   // CompanyID is the unique identifier for the company associated with the condidat.
	CompanyName      string    `json:"companyName"` // CompanyName is the name of the company associated with the condidat.
	LevelOfEducation string    `json:"educationLevel"`
	University       string    `json:"university"`
	Status   bool     `json:"status"`
	RoleName  string    `json:"role_name"`  
} //@name CondidatDetails

// @Description	Signin represents the information required for signing in candidat.
type Signin struct {
	Email    string `json:"email" binding:"required,email,max=255"`     // Email is the email address of the user. It is required, should be a valid email, and maximum length is 255 characters.
	Password string `json:"password" binding:"required,min=10,max=255"` // Password is the user's password. It is required, and its length should be between 10 and 255 characters.
} //@name Signin



// @Description	LoggedInResponse represents the response structure after successful login.
type LoggedInResponse struct {
	AccessToken string   `json:"accessToken"` // AccessToken is the token obtained after successful login for authentication purposes.
	Candidat        LoggedIn `json:"Candidat"`        // User is the structure containing details of the logged-in user.
} //@name LoggedInResponse

// @Description	UpdateCandidate represents the input structure for creating a new condidat.
type UpdateCandidate struct {

	Status   bool     `json:"status"`
} //@name UpdateCandidate



// @Description	LoggedIn represents the candidat details after successful login.
type LoggedIn struct {
	ID        uuid.UUID `json:"ID"`            // ID is the unique identifier for the candidat.
	Name           string    `json:"name"`           // Name is the name of the user..
	Email     string         `json:"email"`         // Email is the email address of the candidat.
	CompanyID uuid.UUID `json:"workCompanyId"` // CompanyID is the unique identifier for the candidat company.
} //@name LoggedIn


// @Description	GenderPercentagesResponse represents porcentage of gender
type AcceptancePercentagesResponse struct {
    Acceptance_Percentage   float64 `json:"Acceptance_Percentage"`   // MalePercentage est le pourcentage d'hommes dans la base de données.
    Refused_Percentage float64 `json:"Refused_Percentage"` // FemalePercentage est le pourcentage de femmes dans la base de données.
}
//@name GenderPercentagesResponse


// @Description	levelPercentagesResponse represents porcentage of gender
type levelPercentagesResponse struct {
Bachelor   float64 `json:"bachelor"`   // MalePercentage est le pourcentage d'hommes dans la base de données.
Master float64 `json:"master"`
Other   float64    `json:"other"`          // FemalePercentage est le pourcentage de femmes dans la base de données.
}
//@name levelPercentagesResponse