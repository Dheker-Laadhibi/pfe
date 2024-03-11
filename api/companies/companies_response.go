package companies

import (
	"time"

	"github.com/google/uuid"
)

// @Description	CompanyIn represents the input structure for creating a new company.
type CompanyIn struct {
	Name string `json:"name" binding:"required,min=3,max=30"` // Name is the name of the company. It is required and should be between 3 and 30 characters.
} //@name CompanyIn

// @Description	CompaniesPagination represents the paginated list of companies.
type CompaniesPagination struct {
	Items      []CompaniesTable `json:"items"`      // Items is a slice containing individual company details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of companies in the entire list.
} //@name CompaniesPagination

// @Description	CompaniesTable represents a single company entry in a table.
type CompaniesTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the company.
	Name      string    `json:"name"`      // Name is the name of the company.
	Email     string    `json:"email"`     // Email is the email address associated with the company.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt is the timestamp indicating when the company entry was created.
} //@name CompaniesTable

// @Description	CompaniesDetails represents detailed information about a specific company.
type CompaniesDetails struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the company.
	Name      string    `json:"name"`      // Name is the name of the company.
	Email     string    `json:"email"`     // Email is the email address associated with the company.
	Website   string    `json:"website"`   // Website is the website URL of the company.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt is the timestamp indicating when the company entry was created.
} //@name CompaniesDetails
