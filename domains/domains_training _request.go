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
    dheker last update 00:03 14 mars
*/

package domains

import (
	"time"

	"github.com/google/uuid"
)

// TrainingRequest represents information about training request demanded by an employee  in the system.
type TrainingRequest struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the training Request
	TrainingTitle   string    `gorm:"column:training_title; not null"`             // Object   of the training Request
	Description     string    `gorm:"column:description; not null"`                // Description of   the training Request
	Reason          string    `gorm:"column:reason; not null;"`                    // why to apply on
	RequestDate     time.Time `gorm:"column:request_date; not null;"`
	DecisionCompany string      `gorm:"column:decision_company;"`
	UserID          uuid.UUID `gorm:"column:user_id;"` // User ID associated with the training Request

}
