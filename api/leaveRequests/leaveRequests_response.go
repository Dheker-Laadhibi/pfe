package leaveRequests

import (
	"time"

	"github.com/google/uuid"
)

// @Description	LeaveIn represents the input structure for updating the "seen" status of LeaveRequests.
type LeaveRequestIn struct {
	Status string `json:"status" binding:"required"` // Seen is a boolean indicating whether the LeaveRequests has been approved or not. It is required.
} //@name LeaveIn

// @Description	LeaveCount represents the count of LeaveRequests.
type LeaveRequestCount struct {
	Count uint `json:"count"` // Count is the number of LeaveRequests.
} //@name LeaveCount

// @Description	LeaveDetails represents detailed information about a specific LeaveRequests.
type LeaveRequestDetails struct {
	ID        uuid.UUID `json:"id"`                                         //ID is the unique identifier for the LeaveRequests.
	StartDate string    `json:"start_date"`                                 // StartDate of the LeaveRequests
	EndDate   string    `json:"end_date"`                                   // EndDate of the LeaveRequests
	Type      string    `gorm:"column:leave_type; not null; default:false"` //  Type of LeaveRequests (e.g., annual, sick, maternity, unpaid)
	Status    string    `json:"status"`                                     // Status of the LeaveRequests request (pending, approved, rejected)
	Reason    string    `json:"reason"`                                     // Reason of the LeaveRequests
	UserID    uuid.UUID `gorm:"column:user_id;"`                            // User ID associated with the Leave requests
	CreatedAt time.Time `json:"createdAt"`                                  // CreatedAt is the timestamp indicating when the LeaveRequests was created.
} //@name LeaveDetails

// @Description	LeaveDemande represents detailed information about a specific LeaveDemande.
type LeaveRequestDemande struct {
	StartDate string `json:"start_date"`                                 // StartDate of the LeaveRequests
	EndDate   string `json:"end_date"`                                   // EndDate of the LeaveRequests
	Type      string `gorm:"column:leave_type; not null; default:false"` //  Type of LeaveRequests (e.g., annual, sick, maternity, unpaid)
	Reason    string `json:"reason"`                                     // Reason of the LeaveRequests
} //@name LeaveDemande

// @Description	LeavePagination represents the paginated list of LeaveRequests.
type LeavePagination struct {
	Items      []LeaveRequestDetails `json:"items"`      // Items is a slice containing individual LeaveRequest details.
	Page       uint                  `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint                  `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint                  `json:"totalCount"` // TotalCount is the total number of LeaveRequests in the entire list.
} //@name LeavePagination
