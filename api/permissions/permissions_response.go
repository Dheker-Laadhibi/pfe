package interns

import (


	"github.com/google/uuid"
)
// @Description	permissionIn represents the input structure for creating a new intern.
type permissionIn struct {
	FeatureName     string    `json:"feature_name"`           
	CreatePerm       bool   `json:"create_perm"` // Firstname is the first name of the intern. It is required and should be between 3 and 30 characters.
	ReadPerm         bool  `json:"read_perm"`  // Lastname is the last name of the intern. It is required and should be between 3 and 35 characters.
	UpdatePerm      bool `json:"update_perm"`    // Email is the email address of the intern. It is required, should be a valid email, and maximum length is 255 characters.
	DeletePerm      bool    `json:"delete_perm"` // intern's education level        
} //@name permissionIn

// @Description	PermissionPagination represents the paginated list of interns.
type PermissionsPagination struct {
	Items      []PermissionTable  `json:"items"`      // Items is a slice containing individual intern details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint            `json:"totalCount"` // TotalCount is the total number of interns in the entire list.
} //@name PermissionPagination

// @Description	PermissionTable represents a single intern entry in a table.
type PermissionTable struct {
	FeatureName     string    `json:"feature_name"` 
	ID           uuid.UUID `json:"id"`        // ID is the unique identifier for the intern.
	CreatePerm   bool `json:"create_perm"` // Firstname is the first name of the intern. It is required and should be between 3 and 30 characters.
	ReadPerm     bool `json:"read_perm"`  // Lastname is the last name of the intern. It is required and should be between 3 and 35 characters.
	UpdatePerm   bool `json:"update_perm"`    // Email is the email address of the intern. It is required, should be a valid email, and maximum length is 255 characters.
	DeletePerm   bool    `json:"delete_perm"` // intern's education level 
} //@name PermissionTable

// @Description	PermissionList represents a simplified version of the intern for listing purposes.
type PermissionList struct {
	ID   uuid.UUID `json:"id"`   // ID is the unique identifier for the intern.
	FeatureName     string    `json:"feature_name"`           
 
} //@name PermissionList

// @Description	InternsCount represents the count of interns.
type InternsCount struct {
	Count uint `json:"count"` // Count is the number of interns.
} //@name InternsCount

// @Description	PermissionsDetails represents detailed information about a specific intern.
type PermissionsDetails struct {
	ID               uuid.UUID `json:"id"`        // ID is the unique identifier for the intern.
	CreatePerm bool `json:"create_perm"` // Firstname is the first name of the intern. It is required and should be between 3 and 30 characters.
	ReadPerm  bool `json:"read_perm"`  // Lastname is the last name of the intern. It is required and should be between 3 and 35 characters.
	UpdatePerm     bool `json:"update_perm"`    // Email is the email address of the intern. It is required, should be a valid email, and maximum length is 255 characters.
	DeletePerm bool    `json:"delete_perm"` // intern's education level  
} //@name PermissionsDetails
