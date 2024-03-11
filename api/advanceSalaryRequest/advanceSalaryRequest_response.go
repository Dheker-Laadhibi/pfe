package advanceSalaryRequest

import (
	"time"

	"github.com/google/uuid"
)

// @Description	AdvanceSalaryRequestIn represents the input structure for updating the "seen" status of AdvanceSalaryRequests.
type AdvanceSalaryRequestIn struct {
	Status string `json:"status" binding:"required"` // Seen is a boolean indicating whether the advanceSalaryRequests has been approved or not. It is required.
} //@name AdvanceSalaryRequestIn

// @Description	LeaveCount represents the count of advanceSalaryRequest.
type AdvanceSalaryRequestCount struct {
	Count uint `json:"count"` // Count is the number of advanceSalaryRequest.
} //@name AdvanceSalaryRequestCount

// @Description	AdvanceSalaryRequestDetails represents detailed information about a specific advanceSalaryRequest.
type AdvanceSalaryRequestDetails struct {
	ID        uuid.UUID `json:"id"`                      //ID is the unique identifier for the advanceSalaryRequest.
	Amount    float64   `gorm:"column:amount; not null"` // The amount of the advance.
	Status    string    `json:"status"`                  // Status of the advanceSalaryRequest request.(pending, approved, rejected)
	Reason    string    `json:"reason"`                  // Reason of the advanceSalaryRequest.
	CreatedAt time.Time `json:"createdAt"`               // CreatedAt is the timestamp indicating when the advanceSalaryRequest was created.
} //@name AdvanceSalaryRequestDetails

// @Description	AdvanceSalaryRequestDemande represents detailed information about a specific AdvanceSalaryRequestDemande.
type AdvanceSalaryRequestDemande struct {
	Amount float64 `gorm:"column:amount; not null"` // The amount of the advance
	Reason string  `json:"reason"`                  // Reason of the advanceSalaryRequests
} //@name AdvanceSalaryRequestDemande

// @Description	AdvanceSalaryRequestPagination represents the paginated list of AdvanceSalaryRequests.
type AdvanceSalaryRequestPagination struct {
	Items      []AdvanceSalaryRequestDetails `json:"items"`      // Items is a slice containing individual advanceSalaryRequest details.
	Page       uint                          `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint                          `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint                          `json:"totalCount"` // TotalCount is the total number of advanceSalaryRequests in the entire list.
} //@name AdvanceSalaryRequestPagination
