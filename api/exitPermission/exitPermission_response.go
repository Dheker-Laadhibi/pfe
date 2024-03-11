package exitPermission

import (
	"time"

	"github.com/google/uuid"
)

// @Description	ExitPermissionIn represents the input structure for updating the status of exitPermission.
type ExitPermissionIn struct {
	Status string `gorm:"column:status; not null; default:pending"` // Status of the exitPermission request (e.g., pending, approved, rejected).
}

// @Description	ExitPermissionCount represents the count of exitPermission.
type ExitPermissionCount struct {
	Count uint `json:"count"` // Count is the number of exitPermission.
} //@name ExitPermissionCount

// @Description	ExitPermissionDemande represents a demande information about a specific .
type ExitPermissionDemande struct {
	Reason     string `gorm:"column:reason; not null"`      // Reason or explanation for requesting the exit permission.
	StartDate  string `gorm:"column:start_date; not null"`  // The start date
	ReturnDate string `gorm:"column:return_date; not null"` // The return date
	Type       string `gorm:"column:type; not null"`        // The type of the exit permission
} //@name ExitPermissionDemande

// @Description	ExitPermissionDetails represents detailed information about a specific exitPermission.
type ExitPermissionDetails struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the exitPermission
	Reason    string    `gorm:"column:type; not null"`                       // Reason or explanation for requesting the exit permission.
	Status    string    `json:"status" binding:"required"`                   // Status of the exitPermission request (e.g., pending, approved, rejected).
	UserID    uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the exitPermission
	CreatedAt time.Time `json:"createdAt"`                                   // CreatedAt is the timestamp indicating when the exitPermission was created.
} //@name ExitPermissionDetails

// @Description	ExitPermissionPagination represents the paginated list of ExitPermission.
type ExitPermissionPagination struct {
	Items      []ExitPermissionDetails `json:"items"`      // Items is a slice containing individual ExitPermission details.
	Page       uint                    `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint                    `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint                    `json:"totalCount"` // TotalCount is the total number of ExitPermission in the entire list.
} //@name ExitPermissionPagination
