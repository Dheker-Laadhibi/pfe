/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Leave requests: Represents information about Leave request in the system.
		- ID (uuid.UUID): Unique identifier for the Leave request.
		- StartDate (string): Type of the Leave request.
		- EndDate (string): The end Date of the Leave request.
		- LeaveType (string): Type of Leave request (e.g., annual, sick, maternity, unpaid).
		- Status (string): Status of the leave request (pending, approved, rejected).
		- Reason (string): Reason of the Leave request.
		- UserID (uuid.UUID): User ID associated with the Leave request.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling Leave requests in the project.

	Note:
	- The Leave requests structure represents information about Leave requests in the system.

	Last update :
	24/02/2024 22:22

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Leave requests represents information about the leave requests in the system.
type LeaveRequests struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the Leave requests
	StartDate string    `gorm:"column:start_date; not null"`                 // StartDate of the Leave requests
	EndDate   string    `gorm:"column:end_date; not null;"`                  // EndDate of the Leave requests
	Type      string    `gorm:"column:leave_type; not null; default:false"`  //  Type of the Leave requests (e.g., annual, sick, maternity, unpaid)
	Status    string    `gorm:"column:status; not null; default:pending"`    // Status of the Leave requests request (pending, approved, rejected)
	Reason    string    `gorm:"column:reason; not null; default:false"`      // Reason of the Leave requests
	UserID    uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the Leave requests
	gorm.Model
}
