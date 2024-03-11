/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Questions: Represents the roles defined in the system.
		- ID (uuid.UUID): Unique identifier for the question
		- Question (string): The text of the question
		- CorrectAnswer (string): The correct answer to the question
		- Options (pq.StringArray): The options of the question
        - AssociatedTechnology (string): Associated technology or subject for the question
		- CompanyID (uuid.UUID): ID of the company associated with the question


	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling Questions in the project.

	Note:
	- The Questions structure represents the Questions defined in the system.

	Last update :
	03/03/2024 10:22

*/

package domains

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Questions represents the Questions available for the tests.
type Questions struct {
	ID                   uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the question
	Question             string         `gorm:"column:question; not null;"`                  // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer; not null;"`            // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`      // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology; not null;"`     // Associated technology or subject for the question
	CompanyID            uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company associated with the question
	gorm.Model
}
