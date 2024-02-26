/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Notifications: Represents information about notifications in the system.
		- ID (uuid.UUID): Unique identifier for the notification.
		- Type (string): Type of the notification.
		- Content (string): Content of the notification.
		- Seen (bool): Indicates whether the notification has been seen or not (default: false).
		- UserID (uuid.UUID): User ID associated with the notification.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling notifications in the project.

	Note:
	- The Notifications structure represents information about notifications in the system.

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notifications represents information about notifications in the system.
type Notifications struct {
	ID      uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the notification
	Type    string    `gorm:"column:type; not null"`                       // Type of the notification
	Content string    `gorm:"column:content; not null;"`                   // Content of the notification
	Seen    bool      `gorm:"column:seen; not null; default:false"`        // Indicates whether the notification has been seen or not (default: false)
	UserID  uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the notification
	gorm.Model
}
