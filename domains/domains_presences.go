/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Presences: Represents information about presences in the system.
		- ID (uuid.UUID): Unique identifier for the presence.
		- Type (string): Type of the presence.
		- Content (string): Content of the presence.
		- Seen (bool): Indicates whether the presence has been seen or not (default: false).
		- UserID (uuid.UUID): User ID associated with the presence.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling presences in the project.

	Note:
	- The Presences structure represents information about presences in the system.

*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Presences represents information about presences in the system.
type Presences struct {
	ID      uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the presence
	Matricule    uint    `gorm:"column:matricule; not null"`                       // Type of the presence
	Check  time.Time    `gorm:"column:check; not null;"`                   // Content of the presence
	UserID  uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the presence
	gorm.Model
}	
