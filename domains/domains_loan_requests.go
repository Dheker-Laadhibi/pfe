/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- loan requests: Represents information about loan requests in the system.
		- ID (uuid.UUID): Unique identifier for the loan requests.
		- LoanAmount(float): The amount of the loan request.
		- LoanDuration(string): The duration of the loan request.
		- InterestRate(float): The interest rate of the loan request.
		- ReasonForLoan(string): Reason or explanation for requesting the loan request.
		- Status(string):  Status of the loan requests request (e.g., pending, approved, rejected).
		- PathDocument(string): The path document for the loan request.
		- UserID (uuid.UUID): User ID associated with the loan requests.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling loan requests in the project.

	Note:
	- The loan requests structure represents information about loan requests in the system.

	Last update :
	24/02/2024 13:20

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// loan requests represents information about loan requests in the system.
type LoanRequests struct {
	ID            uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the loan request
	LoanAmount    float64   `gorm:"column:loan_amount; not null"`                // The amount of the loan request
	LoanDuration  string    `gorm:"column:loan_duration; not null"`              // The duration of the loan request
	InterestRate  float64   `gorm:"column:interest_rate; not null"`              // The interest rate of the loan request
	ReasonForLoan string    `gorm:"column:reason_for_loan; not null"`            // Reason or explanation for requesting the loan.
	Status        string    `gorm:"column:status; not null; default:pending"`    // Status of the loan requests request (e.g., pending, approved, rejected).
	PathDocument  string    `gorm:"column:path_document; not null"`              // The path document for the loan request
	CompanyID     uuid.UUID `json:"companyID" binding:"required"`                // CompanyID is the unique identifier for the company associated with the loan request. It is required.
	UserID        uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the loan request
	gorm.Model
}
