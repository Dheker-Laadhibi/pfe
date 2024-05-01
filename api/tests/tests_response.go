package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @Description	TestsIn represents the input structure for creating a new test.
type TestsIn struct {
	Title        string         `gorm:"column:title; not null;"`                     // The title of the test
	Specialty    string         `gorm:"column:specialty; not null;"`                 // The specialty of the test
	Technologies pq.StringArray `gorm:"column:technologies; type:text[]; not null;"` // The technologies assigned to the test
} //@name TestsIn

// @Description	AssignTestIn represents the input structure for creating a new test.
type AssignTestIn struct {
	TestID uuid.UUID `gorm:"column:tests_id; type:uuid; primaryKey;"` // The testID associated with the candidat
} //@name AssignTestIn

// @Description	AnswersPagination represents the paginated list of tests.
type AnswersPagination struct {
	Items      []CandidatAnswer `json:"items"`      // Items is a slice containing individual Answer details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of tests in the entire list.
} //@name AnswersPagination

// @Description	CandidatAnswer represents the input structure for updating the CandidatAnswer.
type CandidatAnswer struct {
	QuestionID     uuid.UUID `gorm:"column:question_id; primaryKey; type:uuid; not null;"` // The questionID associated with the question
	Question       string    `gorm:"column:question; not null;"`                           // The text of the question
	CorrectAnswer  string    `gorm:"column:correct_answer; not null;"`                     // The correct answer to the question
	CandidatAnswer string    `gorm:"column:candidat_answer; not null;"`                    // The candidat answer to the question
} //@name CandidatAnswer

// @Description	CandidatAnswerIn represents the input structure for updating the CandidatAnswer.
type CandidatAnswerIn struct {
	CandidatAnswer string `gorm:"column:candidat_answer; not null;"` // The candidat answer to the question
} //@name CandidatAnswerIn

// @Description	TestsQuestionsIn represents the input structure for creating a new test.
type TestsQuestionsIn struct {
	QuestionID           uuid.UUID      `gorm:"column:question_id; primaryKey; type:uuid; not null;"` // The questionID associated with the question
	Question             string         `gorm:"column:question; not null;"`                           // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer; not null;"`                     // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`               // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology; not null;"`              // Associated technology or subject for the question
	CandidatID           uuid.UUID      `gorm:"column:candidat_id; type:uuid; not null;"`             // ID of the candidat associated with the test
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
	Title        string         `gorm:"column:title; not null;"`     // The title of the test
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
	ID                   uuid.UUID      `json:"id"`                                                   // ID is the unique identifier for the test.
	QuestionID           uuid.UUID      `gorm:"column:question_id; primaryKey; type:uuid; not null;"` // The questionID associated with the question
	Question             string         `gorm:"column:question; not null;"`                           // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer; not null;"`                     // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`               // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology; not null;"`              // Associated technology or subject for the question
	CandidatID           uuid.UUID      `gorm:"column:candidat_id; type:uuid; not null;"`             // ID of the candidat associated with the test
} //@name TestsQuestionsTable

// @Description	ScoresPagination represents the paginated list of Scores.
type ScoresPagination struct {
	Items      []ScorsTable `json:"items"`      // Items is a slice containing individual test details.
	Page       uint         `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint         `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint         `json:"totalCount"` // TotalCount is the total number of tests in the entire list.
} //@name ScosPagination

// @Description	QuestionsPagination represents the paginated list of questions by testID.
type QuestionsPagination struct {
	Items      []QuestionsTable `json:"items"`      // Items is a slice containing individual test details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of tests in the entire list.
} //@name TestsPagination

// @Description	TestsTable represents a single question entry in a table.
type QuestionsTable struct {
	QuestionID           uuid.UUID      `gorm:"column:question_id; primaryKey; type:uuid; not null;"` // The questionID associated with the question
	Question             string         `gorm:"column:question; not null;"`                           // The text of the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`               // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology; not null;"`              // Associated technology or subject for the question
	CreatedAt            time.Time      `json:"createdAt"`                                            // CreatedAt is the timestamp indicating when the test entry was created.
} //@name QuestionsTable

// @Description	ScoresTable represents a single Scores entry in a table.
type ScorsTable struct {
	IdTest     uuid.UUID `json:"id_test"`                      // ID is the unique identifier for the test.
	IdCandidat uuid.UUID `json:"id_candidat"`                  // ID is the unique identifier for the candidat.
	Title      string    `gorm:"column:title; not null;"`      // The title of the test
	Firstname  string    `gorm:"column:first_name; not null;"` // The user's first name
	Lastname   string    `gorm:"column:last_name; not null;"`  // The user's last name
	Score      uint      `gorm:"column:score; not null;"`      // The score associated with the test and the candidat
} //@name ScoresTable

// @Description	TestsCandidatsTable represents a single test entry in a table.
type TestsCandidatsList struct {
	ID        uuid.UUID      `json:"id"`                             // ID is the unique identifier for the test.
	Candidats pq.StringArray `gorm:"column:cuestions; type:text[];"` // The candidats that belongs to the test.
} //@name TestsCandidatsTable

// @Description	TestsList represents a simplified version of the test for listing purposes.
type TestsList struct {
	ID           uuid.UUID      `json:"id"`                          // ID is the unique identifier for the test.
	Specialty    string         `gorm:"column:specialty; not null;"` // The specialty of the test
	Title        string         `gorm:"column:title; not null;"`     // The title of the test
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
