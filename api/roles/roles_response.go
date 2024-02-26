package roles

import (
	"time"

	"github.com/google/uuid"
)

// @Description	RolesIn represents the input structure for creating a new role.
type RolesIn struct {
	Name string `json:"name" binding:"required,min=2,max=40"` // Name is the name of the role. It is required and should be between 2 and 40 characters.
} //@name RolesIn

// @Description	RolesPagination represents the paginated list of roles.
type RolesPagination struct {
	Items      []RolesTable `json:"items"`      // Items is a slice containing individual role details.
	Page       uint         `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint         `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint         `json:"totalCount"` // TotalCount is the total number of roles in the entire list.
} //@name RolesPagination

// @Description	RolesTable represents a single role entry in a table.
type RolesTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the role.
	Name      string    `json:"name"`      // Name is the name of the role.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt is the timestamp indicating when the role entry was created.
} //@name RolesTable

// @Description	RolesList represents a simplified version of the role for listing purposes.
type RolesList struct {
	ID   uuid.UUID `json:"id"`   // ID is the unique identifier for the role.
	Name string    `json:"name"` // Name is the name of the role.
} //@name RolesList

// @Description	RolesCount represents the count of roles.
type RolesCount struct {
	Count uint `json:"count"` // Count is the number of roles.
} //@name RolesCount

// @Description	RolesDetails represents detailed information about a specific role.
type RolesDetails struct {
	ID          uuid.UUID `json:"id"`          // ID is the unique identifier for the role.
	Name        string    `json:"firstname"`   // Name is the name of the role.
	CompanyID   uuid.UUID `json:"companyID"`   // CompanyID is the unique identifier for the company associated with the role.
	CompanyName string    `json:"companyName"` // CompanyName is the name of the company associated with the role.
	CreatedAt   time.Time `json:"createdAt"`   // CreatedAt is the timestamp indicating when the role entry was created.
} //@name RolesDetails
