/*
	Package domains provides the data structures representing entities in the project.

	Structures:
	- Tests: Represents the test information in the system.
		- ID (uuid.UUID): Unique identifier for the test.
		- Specialty (string): The specialty of the test.
		- Technologies (pq.StringArray): The technologies assigned to the test.
		- Questions ([]Questions): The questions assigned to the test.
		- Candidats ([]Condidats): The candidats who belongs to that test.
		- CompanyID (uuid.UUID): ID of the company to which the test belongs.

	- TestQuestions: Represents the questions assigned to tests.
		- TestID (uuid.UUID): Test's ID.
		- RandomQuestions (pq.StringArray): Questions's ID.
		- CompanyID (uuid.UUID): ID of the company associated with the test and questions.

	- TestCandidats: Represents the candidats that belongs to the test.
		- TestID (uuid.UUID): Test's ID.
		- CandidatID (uuid.UUID): Candidat's ID.
		- score (uint): Candidat's score.
		- CompanyID (uuid.UUID): ID of the company associated with the candidat and test.

	Dependencies:
	- "errors": Standard Go package for errors handling.
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.
	- "time": Standard Go package for handling time.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling test information in the project.

	Note:
	- The Tests structure represents the test information in the system.

	Last update :
	04/03/2024 10:22

*/

package domains

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Tests represents the test information in the system.
type Tests struct {
	ID           uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the test
	Specialty    string         `gorm:"column:specialty; not null;"`                 // The specialty of the test
	Technologies pq.StringArray `gorm:"column:technologies; type:text[]; not null;"` // The technologies assigned to the test
	Questions    []Questions    `gorm:"many2many:testQuestions"`                     // The questions assigned to the test
	Candidats    []Condidats    `gorm:"many2many:testCandidats"`                     // The candidats who belongs to that test
	CompanyID    uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the test belongs
	gorm.Model
}

// TestQuestions represents the questions assigned to test.
type TestQuestions struct {
	TestID          uuid.UUID      `gorm:"column:tests_id; primaryKey;"`            // The testID associated with the questions
	RandomQuestions pq.StringArray `gorm:"column:random_questions; type:text[];"`   // The questions associated with the test
	CompanyID       uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"` // ID of the company associated with the question and the test
}

// TestCandidats represents the candidats assigned to tests.
type TestCandidats struct {
	TestID     uuid.UUID `gorm:"column:tests_id; type:uuid; primaryKey;"`    // The testID associated with the candidat
	CandidatID uuid.UUID `gorm:"column:candidat_id; type:uuid; primaryKey;"` // The candidatID associated with the test
	Score      uint      `gorm:"column:score; not null;"`                    // The score associated with the test and the candidat
	CompanyID  uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`    // ID of the company associated with the test
}
