package condidats

import (
	"time"

	"github.com/google/uuid"
)

// @Description	RolesIn represents the input structure for creating a new role.
type CondidatIn struct {
	Firstname string    `json:"firstName" binding:"required,min=3,max=30"`  // Firstname is the first name of the user. It is required and should be between 3 and 30 characters.
	Lastname  string    `json:"lastName" binding:"required,min=3,max=35"`   // Lastname is the last name of the user. It is required and should be between 3 and 35 characters.
	Email     string    `json:"email" binding:"required,email,max=255"`     // Email is the email address of the user. It is required, should be a valid email, and maximum length is 255 characters.
	Password  string    `json:"password" binding:"required,min=10,max=255"` // Password is the user's password. It is required, and its length should be between 10 and 255 characters.
	CompanyID uuid.UUID `json:"companyID" binding:"required"`               // CompanyID is the unique identifier for the company associated with the user. It is required.
} //@name RolesIn

// @Description	CondidtasPagination represents the paginated list of Condidats.
type CondidtasPagination struct {
	Items      []CondidatsTable `json:"items"`      // Items is a slice containing individual role details.
	Page       uint         `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint         `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint         `json:"totalCount"` // TotalCount is the total number of roles in the entire list.
} //@name CondidtasPagination

// @Description	RolesTable represents a single role entry in a table.
type CondidatsTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the  condidat.
	Firstname string    `json:"firstname"` // Firstname is the first name of the condidat.
	Lastname  string    `json:"lastname"`  // Lastname is the last name of the condidat.
	Email     string    `json:"email"`     // Email is the email address of the condidat.
} //@name CondidatsTable

// @Description	CondidatsList represents a simplified version of the Condidats for listing purposes.
type CondidatsList struct {
	ID   uuid.UUID `json:"id"`   // ID is the unique identifier for the role.
	Firstname string    `json:"firstname"` // Name is the name of the role.
	Lastname string    `json:"lastname"` // Name is the name of the role.

} //@name CondidatsList

// @Description	RolesCount represents the count of roles.
type CondidatsCount struct {
	Count uint `json:"count"` // Count is the number of roles.
} //@name RolesCount

// @Description	RolesDetails represents detailed information about a specific role.
type CondidatDetails struct {
	ID          uuid.UUID `json:"id"`          // ID is the unique identifier for the role.
	Firstname        string    `json:"firstname"`   // Name is the name of the role.
	Lastname string    `json:"lastname"`  // last name of condidat 
	CompanyID   uuid.UUID `json:"companyID"`   // CompanyID is the unique identifier for the company associated with the role.
	CompanyName string    `json:"companyName"` // CompanyName is the name of the company associated with the role.
	LevelOfEducation  string    `json:"educationLevel"`        
	University       string      `json:"university"`
	CreatedAt   time.Time `json:"createdAt"`   // CreatedAt is the timestamp indicating when the role entry was created.
} //@name RolesDetails
