package projects

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @Description	ProjectIn represents the input structure for creating a new project.
type ProjectIn struct {
	Code         string         `json:"code" binding:"required,min=3,max=30"`        // code is the code of pfe project
	Projectname  string         `json:"projectname" binding:"required,min=3,max=35"` // projectname is the  name of the project. It is required and should be between 3 and 35 characters.
	Description  string         `json:"description" binding:"required,min=3,max=80"`
	Technologies pq.StringArray `json:"technologies" binding:"required"` // technologies required to develop the project
     ExpDate       string       `json:"exp_date" binding:"required"`  
	CompanyID    uuid.UUID      `json:"companyID" binding:"required"`    // CompanyID is the unique identifier for the company associated with the project. It is required.
} //@name ProjectIn

// @Description	ProjectPagination represents the paginated list of projects.
type projectsPagination struct {
	Items      []ProjectTable `json:"items"`      // Items is a slice containing individual project details.
	Page       uint           `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint           `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint           `json:"totalCount"` // TotalCount is the total number of project in the entire list.
} //@name UsersPagination

// @Description	UsersTable represents a single project entry in a table.
type ProjectTable struct {
	ID           uuid.UUID      `json:"id"`           // ID is the unique identifier for the project.
	Code         string         `json:"code"`         // code is the code of pfe project
	Projectname  string         `json:"projectname"`  // projectname is the  name of the project. It is required .
	Technologies pq.StringArray `json:"technologies"` // technologies required to develop the project
	CompanyID    uuid.UUID      `json:"companyID"`    // CompanyID is the unique identifier for the company associated with the project. It is required.
	ExpDate      time.Time      `json:"expdate"`      // expdate is the timestamp indicating when the project entry will ends.
} //@name UsersTable

// @Description	ProjectList represents a simplified version of the project for listing purposes.
type ProjectsList struct {
	ID          uuid.UUID `json:"id"`          // ID is the unique identifier for the project.
	Projectname string    `json:"projectname"` //  projectname is the  name of the project.
	Code        string    `json:"code"`        // code is the code of pfe project
} //@name ProjectList

// @Description	ProjectsCount represents the count of projects.
type ProjectsCount struct {
	Count uint `json:"count"` // Count is the number of projects.
} //@name ProjectsCount

// @Description	projectsDetails represents detailed information about a specific project.
type projectsDetails struct {
	ID           uuid.UUID      `json:"id"`           // ID is the unique identifier for the project.
	Code         string         `json:"code"`         // code is the code of pfe project
	Projectname  string         `json:"projectname"`  // projectname is the  name of the project. It is required .
	Technologies pq.StringArray `json:"technologies"` // technologies required to develop the project
	CompanyID    uuid.UUID      `json:"companyID"`    // CompanyID is the unique identifier for the company associated with the project. It is required.
	ExpDate      time.Time      `json:"expdate"`      // expdate is the timestamp indicating when the project entry will ends.
} //@name projectsDetails
