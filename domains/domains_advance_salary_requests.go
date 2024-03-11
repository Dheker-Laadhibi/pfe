/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- AdvanceSalaryRequest: Represents information about advance salary requests in the system.
		- ID (uuid.UUID): Unique identifier for the advance salary requests.
		- Amount (float): The amount of the advance
		- Status (string): Status of the advance salary requests (pending, approved, rejected).
		- Reason (string): Reason of the advance salary requests
		- UserID (uuid.UUID): User ID associated with the advance salary requests.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling advance salary requests in the project.

	Note:
	- The AdvanceSalaryRequests structure represents information about advance salary requests in the system.

	Last update :
	24/02/2024 22:22

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// advance salary requests represents information about leave_requests in the system.
type AdvanceSalaryRequests struct {
	ID     uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the advance salary request
	Amount float64   `gorm:"column:amount; not null"`                     // The amount of the advance
	Reason string    `gorm:"column:reason; not null; default:false"`      // Reason of the advance salary request
	Status string    `gorm:"column:status; not null; default:pending"`    // Status of the advance salary requests (pending, approved, rejected)
	UserID uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the advance salary request
	gorm.Model
}
