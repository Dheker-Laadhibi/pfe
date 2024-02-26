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

// Presences represents information about mission_orders  in the system.
type MissionOrders struct {
	ID           uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the missionOrders
	Object       string    `gorm:"column:matricule; not null"`                  // Object   of the missionOrders
	Description  string    `gorm:"column:descriptino; not null"`                // Description of   the missionOrders
	StartDate    time.Time `gorm:"column:start_date; not null;"`                // StartDate of the missionOrders
	EndDate      time.Time `gorm:"column:end_date; not null;"`                  // EndDte    of the missionOrders
	AdressClient string    `gorm:"column:adress_client; not null;"`             // AdressClient  of the missionOrders
	Transport    string    `gorm:"column:transport; not null;"`                 // Transport of the missionOrders
	UserID       uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the missionOrders
	gorm.Model
}


