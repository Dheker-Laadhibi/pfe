package interns

import (
	"time"

	"github.com/google/uuid"
)



// @Description	InternsIn represents the input structure for creating a new intern.
type InternsIn struct {
	Firstname string `json:"firstName"` // Firstname is the first name of the intern. It is required and should be between 3 and 30 characters.
	Lastname  string `json:"lastName"`  // Lastname is the last name of the intern. It is required and should be between 3 and 35 characters.
	Email     string `json:"email"`    // Email is the email address of the intern. It is required, should be a valid email, and maximum length is 255 characters.
	//	Lastname        string    `gorm:"column:last_name; not null;"`
	LevelOfEducation string    `json:"educationLevel"` // intern's education level 
	University       string    `json:"university"` // university
	StartDate        string `json:"start_date"`      // start date of the internship
	EndDate          string  `json:"end_date"` 
	Gender           string  
	Adress           string 
	PhoneNumber      string               
} //@name InternsIn

// @Description	InternsPagination represents the paginated list of interns.
type InternsPagination struct {
	Items      []InternsTable `json:"items"`      // Items is a slice containing individual intern details.
	Page       uint           `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint           `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint           `json:"totalCount"` // TotalCount is the total number of interns in the entire list.
} //@name InternsPagination

// @Description	InternsTable represents a single intern entry in a table.
type InternsTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the intern.
	Firstname string    `json:"firstname"` // Firstname is the first name of the intern.
	Lastname  string    `json:"lastname"`  // Lastname is the last name of the intern.
	Email     string    `json:"email"`     // Email is the email address of the intern.
	LevelOfEducation string   `json:"educationLevel"` // intern's education level 
	University       string    `json:"university"` // university
	StartDate        time.Time  `json:"start_date"`      // start date of the internship
	EndDate          time.Time   `json:"end_date"`  
	Gender           string  
	Adress           string 
	PhoneNumber      string  
} //@name InternsTable

// @Description	InternsList represents a simplified version of the intern for listing purposes.
type InternsList struct {
	ID   uuid.UUID `json:"id"`   // ID is the unique identifier for the intern.
	Name string    `json:"name"` // Name is the full name of the intern.
} //@name InternsList

// @Description	InternsCount represents the count of interns.
type InternsCount struct {
	Count uint `json:"count"` // Count is the number of interns.
} //@name InternsCount

// @Description	InternsDetails represents detailed information about a specific intern.
type InternsDetails struct {
	ID               uuid.UUID `json:"id"`        // ID is the unique identifier for the intern.
	Firstname        string    `json:"firstname"` // Firstname is the first name of the intern.
	Lastname         string    `json:"lastname"`  // Lastname is the last name of the intern.
	Email            string    `json:"email"`     // Email is the email address of the intern.
	LevelOfEducation string    `json:"educationLevel"` // intern's education level 
	University       string    `json:"university"` // university
	StartDate        time.Time `json:"startDate"` // startDate of internship 
	Gender           string  
	Adress           string 
	PhoneNumber      string  
} //@name InternsDetails
