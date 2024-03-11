package questions

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @Description	QuestionsIn represents the input structure for creating a new question.
type QuestionsIn struct {
	Question             string         `gorm:"column:question; not null;"`               // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer;not null;"`          // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`   // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology ; not null;"` // Associated technology or subject for the question
} //@name QuestionsIn

// @Description	QuestionsPagination represents the paginated list of questions.
type QuestionsPagination struct {
	Items      []QuestionsTable `json:"items"`      // Items is a slice containing individual question details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of questions in the entire list.
} //@name QuestionsPagination

// @Description	QuestionsTable represents a single question entry in a table.
type QuestionsTable struct {
	ID                   uuid.UUID      `json:"id"`                                       // ID is the unique identifier for the question.
	Question             string         `gorm:"column:question; not null;"`               // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer;not null;"`          // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`   // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology ; not null;"` // Associated technology or subject for the question
	CreatedAt            time.Time      `json:"createdAt"`                                // CreatedAt is the timestamp indicating when the question entry was created.
} //@name QuestionsTable

// @Description	QuestionsList represents a simplified version of the question for listing purposes.
type QuestionsList struct {
	ID                   uuid.UUID      `json:"id"`                                       // ID is the unique identifier for the question.
	Question             string         `gorm:"column:question; not null;"`               // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer;not null;"`          // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:text[]; not null;"`   // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology ; not null;"` // Associated technology or subject for the question
} //@name QuestionsList

// @Description	QuestionsCount represents the count of questions.
type QuestionsCount struct {
	Count uint `json:"count"` // Count is the number of questions.
} //@name QuestionsCount

// @Description	QuestionsDetails represents detailed information about a specific question.
type QuestionsDetails struct {
	ID                   uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the Question
	Question             string         `gorm:"column:question; not null;"`                  // The text of the question
	CorrectAnswer        string         `gorm:"column:correct_answer;not null;"`             // The correct answer to the question
	Options              pq.StringArray `gorm:"column:options; type:jsonb; not null;"`       // The options of the question
	AssociatedTechnology string         `gorm:"column:associated_technology ; not null;"`    // Associated technology or subject for the question
	CompanyName          string         `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company associated with the question
	CreatedAt            time.Time      `json:"createdAt"`                                   // CreatedAt is the timestamp indicating when the question entry was created.
} //@name QuestionsDetails
