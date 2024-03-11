package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @Description	TestsIn represents the input structure for creating a new test.
type TestsIn struct {
	Specialty    string         `gorm:"column:specialty; not null;"`                 // The specialty of the test
	Technologies pq.StringArray `gorm:"column:technologies; type:text[]; not null;"` // The technologies assigned to the test
} //@name TestsIn

// @Description	AssignTestIn represents the input structure for creating a new test.
type AssignTestIn struct {
	TestID uuid.UUID `gorm:"column:tests_id; type:uuid; primaryKey;"` // The testID associated with the candidat
} //@name AssignTestIn

// @Description	TestsQuestionsIn represents the input structure for creating a new test.
type TestsQuestionsIn struct {
	RandomQuestions pq.StringArray `gorm:"column:random_questions; type:text[]; not null;"`
} //@name TestsQuestionsIn

// @Description	TestsPagination represents the paginated list of tests.
type TestsPagination struct {
	Items      []TestsTable `json:"items"`      // Items is a slice containing individual test details.
	Page       uint         `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint         `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint         `json:"totalCount"` // TotalCount is the total number of tests in the entire list.
} //@name TestsPagination

// @Description	TestsTable represents a single test entry in a table.
type TestsTable struct {
	ID           uuid.UUID      `json:"id"`                          // ID is the unique identifier for the test.
	Specialty    string         `gorm:"column:specialty; not null;"` // The user's first name
	Technologies pq.StringArray `gorm:"column:technologies; type:text[]; not null;"`
	CreatedAt    time.Time      `json:"createdAt"` // CreatedAt is the timestamp indicating when the test entry was created.
} //@name TestsTable

// @Description	TestsQuestionsPagination represents the paginated list of tests.
type TestsQuestionsPagination struct {
	Items      []TestsQuestionsTable `json:"items"`      // Items is a slice containing individual test details.
	Page       uint                  `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint                  `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint                  `json:"totalCount"` // TotalCount is the total number of tests in the entire list.
} //@name TestsQuestionsPagination

// @Description	TestsQuestionsTable represents a single test entry in a table.
type TestsQuestionsTable struct {
	ID              uuid.UUID      `json:"id"`                                    // ID is the unique identifier for the test.
	RandomQuestions pq.StringArray `gorm:"column:random_questions; type:text[];"` // The questions associated with the test
} //@name TestsQuestionsTable

// @Description	TestsCandidatsTable represents a single test entry in a table.
type TestsCandidatsList struct {
	ID        uuid.UUID      `json:"id"`                             // ID is the unique identifier for the test.
	Candidats pq.StringArray `gorm:"column:cuestions; type:text[];"` // The candidats that belongs to the test.
} //@name TestsCandidatsTable

// @Description	TestsList represents a simplified version of the test for listing purposes.
type TestsList struct {
	ID           uuid.UUID      `json:"id"`                          // ID is the unique identifier for the test.
	Specialty    string         `gorm:"column:specialty; not null;"` // The specialty of the test
	Technologies pq.StringArray `gorm:"column:technologies; type:text[]; not null;"`
} //@name TestsList

// @Description	TestsCount represents the count of tests.
type TestsCount struct {
	Count uint `json:"count"` // Count is the number of tests.
} //@name TestsCount

// @Description	TestsDetails represents detailed information about a specific test.
type TestsDetails struct {
	ID           uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the test
	Specialty    string         `gorm:"column:specialty; not null;"`                 // The specialty of the test
	Technologies pq.StringArray `gorm:"column:technologies; type:text[]; not null;"` // The technologies assigned to the test
	CompanyID    uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the test belongs
	CreatedAt    time.Time      `json:"createdAt"`                                   // CreatedAt is the timestamp indicating when the test entry was created.
} //@name TestsDetails
