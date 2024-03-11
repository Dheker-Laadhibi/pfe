package loanRequests

import (
	"time"

	"github.com/google/uuid"
)

// @Description	LoanRequestIn represents the input structure for updating the status of LoanRequest.
type LoanRequestIn struct {
	Status string `gorm:"column:status; not null; default:pending"` // Status of the LoanRequest request (e.g., pending, approved, rejected).
}

// @Description	LoanRequestCount represents the count of LoanRequest.
type LoanRequestCount struct {
	Count uint `json:"count"` // Count is the number of LoanRequest.
} //@name LoanRequestCount

// @Description	LoanRequestDemande represents a demande information about a specific .
type LoanRequestDemande struct {
	LoanAmount    float64 `gorm:"column:loan_amount; not null"`     // The amount of the loan request
	LoanDuration  string  `gorm:"column:loan_duration; not null"`   // The duration of the loan request
	InterestRate  float64 `gorm:"column:interest_rate; not null"`   // The interest rate of the loan request
	ReasonForLoan string  `gorm:"column:reason_for_loan; not null"` // Reason or explanation for requesting the loan.
	PathDocument  string  `gorm:"column:path_document; not null"`   // The path document for the loan request
} //@name LoanRequestDemande

// @Description	LoanRequestDetails represents detailed information about a specific LoanRequest.
type LoanRequestDetails struct {
	ID            uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the LoanRequest
	LoanAmount    float64   `gorm:"column:loan_amount; not null"`                // The amount of the loan request
	LoanDuration  string    `gorm:"column:loan_duration; not null"`              // The duration of the loan request
	InterestRate  float64   `gorm:"column:interest_rate; not null"`              // The interest rate of the loan request
	ReasonForLoan string    `gorm:"column:reason_for_loan; not null"`            // Reason or explanation for requesting the loan.
	Status        string    `json:"status" binding:"required"`                   // Status of the LoanRequest request (e.g., pending, approved, rejected).
	PathDocument  string    `gorm:"column:path_document; not null"`              // The path document for the loan request
	UserID        uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the LoanRequest
	CreatedAt     time.Time `json:"createdAt"`                                   // CreatedAt is the timestamp indicating when the loanRequests was created.
} //@name LoanRequestDetails

// @Description	LoanRequestPagination represents the paginated list of LoanRequest.
type LoanRequestPagination struct {
	Items      []LoanRequestDetails `json:"items"`      // Items is a slice containing individual LoanRequest details.
	Page       uint                 `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint                 `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint                 `json:"totalCount"` // TotalCount is the total number of LoanRequest in the entire list.
} //@name LoanRequestPagination
