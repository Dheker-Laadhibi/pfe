package features

import (
	

	"github.com/google/uuid"
)

// @Description	FeatureIn represents the input structure for creating a new role.
type FeatureIn struct {
	Featurename string `json:"feature_name" binding:"required,min=2,max=40"` // Name is the name of the role. It is required and should be between 2 and 40 characters.
} //@name FeatureIn

// @Description FeaturePagination represents the paginated list of roles.
type FeaturePagination struct {
	Items      []FeatureTable `json:"items"`      // Items is a slice containing individual role details.
	Page       uint           `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint           `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint           `json:"totalCount"` // TotalCount is the total number of roles in the entire list.
} //@name FeaturePagination

// @Description	FeatureTable represents a single role entry in a table.
type FeatureTable struct {
	ID        uuid.UUID        `json:"id"`        // ID is the unique identifier for the role.
	Featurename      string    `json:"feature_name"`      // Name is the name of the role.

} //@name FeatureTable


// @Description	FeatureCount represents the count of roles.
type FeatureCount struct {
	Count uint `json:"count"` // Count is the number of roles.
} //@name FeatureCount

