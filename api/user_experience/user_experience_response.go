package user_experience

import (
	"github.com/google/uuid"
)

// @Description	ExperienceIn represents the input structure for creating a new experiences for one user .
type ExperienceIn struct {
	ProfessionalTraining string `json:"professional_training" binding:"required"` //experiences that wil be added to the specified user
} //@name ExperienceIn

// @Description	ExperiencePagination represents the paginated list of projects.
type ExperiencePagination struct {
	Items      []ExperienceTable `json:"items"`      // Items is a slice containing individual project details.
	Page       uint              `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint              `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint              `json:"totalCount"` // TotalCount is the total number of project in the entire list.
} //@name ExperiencePagination

// @Description	ExperienceTable represents a single project entry in a table.
type ExperienceTable struct {
	ID uuid.UUID `json:"id"` // ID is the unique identifier for the project.
	//ProfessionalTrainings pq.StringArray `json:"professional_trainings" binding:"required"`
	ProfessionalTraining string `json:"professional_training"` // technologies required to develop the project

} //@name ExperienceTable

// @Description	ExperienceList represents a simplified version of the project for listing purposes.
type ExperienceList struct {
	ID                   uuid.UUID `json:"id"`                     // ID is the unique identifier for the project.
	ProfessionalTraining string    `json:"professional_trainings"` // technologies required to develop the project
} //@name ExperienceList

// @Description	ProjectsCount represents the count of projects.
type ExperienceCount struct {
	Count uint `json:"count"` // Count is the number of projects.
} //@name ExperienceCount

// @Description	ExperienceDetails represents detailed information about a specific user exp.
type ExperienceDetails struct {
	ID                   uuid.UUID `json:"id"`                    // ID is the unique identifier for the exp.
	ProfessionalTraining string    `json:"professional_training"` //  the user exp
	CompanyID            uuid.UUID `json:"CompanyID"`
} //@name ExperienceDetails
