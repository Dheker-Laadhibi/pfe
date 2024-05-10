package users

import (
	"time"

	"github.com/google/uuid"
)

// @Description	UsersIn represents the input structure for creating a new user.
type UsersIn struct {
	Firstname string    `json:"firstName" binding:"required,min=3,max=30"`  // Firstname is the first name of the user. It is required and should be between 3 and 30 characters.
	Lastname  string    `json:"lastName" binding:"required,min=3,max=35"`   // Lastname is the last name of the user. It is required and should be between 3 and 35 characters.
	Email     string    `json:"email" binding:"required,email,max=255"`     // Email is the email address of the user. It is required, should be a valid email, and maximum length is 255 characters.
	Password  string    `json:"password" binding:"required,min=10,max=255"` // Password is the user's password. It is required, and its length should be between 10 and 255 characters.
	RoleName  string    `json:"role_name"`                                  // RoleName is the name of the role associated with the user.
	CompanyID uuid.UUID `json:"companyID" binding:"required"`               // CompanyID is the unique identifier for the company associated with the user. It is required.
	Gender      string    `json:"gender"`    // gender 

} //@name UsersIn

// @Description	UsersPagination represents the paginated list of users.
type UsersPagination struct {
	Items      []UsersTable `json:"items"`      // Items is a slice containing individual user details.
	Page       uint         `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint         `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint         `json:"totalCount"` // TotalCount is the total number of users in the entire list.
} //@name UsersPagination

// @Description	UsersTable represents a single user entry in a table.
type UsersTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the user.
	Firstname string    `json:"firstname"` // Firstname is the first name of the user.
	Lastname  string    `json:"lastname"`  // Lastname is the last name of the user.
	Email     string    `json:"email"`     // Email is the email address of the user.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt is the timestamp indicating when the user entry was created.
} //@name UsersTable

// @Description	UsersList represents a simplified version of the user for listing purposes.
type UsersList struct {
	ID   uuid.UUID `json:"id"`   // ID is the unique identifier for the user.
	Name string    `json:"name"` // Name is the full name of the user.
} //@name UsersList

// @Description	UsersCount represents the count of users.
type UsersCount struct {
	Count uint `json:"count"` // Count is the number of users.
} //@name UsersCount

// @Description	UsersDetails represents detailed information about a specific user.
type UsersDetails struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the user.
	Firstname string    `json:"firstname"` // Firstname is the first name of the user.
	Lastname  string    `json:"lastname"`  // Lastname is the last name of the user.
	Email     string    `json:"email"`     // Email is the email address of the user.
	Country   string    `json:"country"`   // Country is the country of residence of the user.
	Status    bool      `json:"status"`    // Status is a boolean indicating the current status of the user.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt is the timestamp indicating when the user entry was created.
} //@name UsersDetails

// @Description	UsersList represents a simplified version of the user for listing purposes.
type AddTrainingUser struct {
	ID uuid.UUID `json:"id"` // ID is the unique identifier for the user.
}

// @Description	GenderPercentagesResponse represents porcentage of gender
type GenderPercentagesResponse struct {
    MalePercentage   float64 `json:"malePercentage"`   // MalePercentage est le pourcentage d'hommes dans la base de données.
    FemalePercentage float64 `json:"femalePercentage"` // FemalePercentage est le pourcentage de femmes dans la base de données.
}
//@name GenderPercentagesResponse