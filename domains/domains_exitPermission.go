/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- ExitPermission: Represents information about exitPermission in the system.
		- ID (uuid.UUID): Unique identifier for the exitPermission.
		- ReleaseDate(Time): Release time of the exit permission
		- StartTime(Time): Start time
		- ReturnTime(Time): Return time
		- Type(string): The type of the exit permission
		- Reason (string): Reason or explanation for requesting the exit permission.
		- Status (string):  Status of the exitPermission request (e.g., pending, approved, rejected).
		- UserID (uuid.UUID): User ID associated with the exitPermission.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling exitPermission in the project.

	Note:
	- The ExitPermission structure represents information about exitPermission in the system.

	Last update :
	24/02/2024 13:20

*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExitPermission represents information about the exit Permissions in the system.
type ExitPermission struct {
	ID          uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the exitPermission
	ReleaseDate time.Time `gorm:"column:release_date; not null"`               // Release time of the exit permission
	StartTime   time.Time `gorm:"column:start_time; not null"`                 // Start time
	ReturnTime  string    `gorm:"column:return_time; not null"`                // Return time
	Type        string    `gorm:"column:type; not null"`                       // The type of the exit permission
	Reason      string    `gorm:"column:reason; not null"`                     // Reason or explanation for requesting the exit permission.
	Status      string    `gorm:"column:status; not null; default:pending"`    // Status of the exitPermission request (e.g., pending, approved, rejected).
	UserID      uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the exitPermission
	gorm.Model
}
