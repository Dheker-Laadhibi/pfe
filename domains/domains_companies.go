/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Companies: Represents information about a company in the system.
		- ID (uuid.UUID): Unique identifier for the company.
		- Name (string): The name of the company.
		- Website (string): The company's website URL.
		- Email (string): The company's email address.
		- Team ([]Users): List of users associated with the company.
		- CreatedByUserID (uuid.UUID): ID of the user who created the company.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- ReadCompanyNameByID(db *gorm.DB, companyID uuid.UUID) (string, error): Retrieves the name of the company based on its ID from the database.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling entities in the project.

	Note:
	- The Companies structure represents information about a company in the system.
	- ReadCompanyNameByID retrieves the name of the company based on its ID from the database.

	Last update :
	01/02/2024 10:22

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Companies represents information about a company in the system.
type Companies struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`                                         // Unique identifier for the company
	Name            string    `gorm:"column:name; not null;"`                                                              // The company's name
	Email           string    `gorm:"column:email;"`                                                                       // The company's email
	Website         string    `gorm:"column:website;"`                                                                     // The company's website
	Team            []Users   `gorm:"foreignKey:company_id; references:id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // List of users associated with the company
	CreatedByUserID uuid.UUID `gorm:"column:created_by_user_id; type:uuid; not null;"`                                     // ID of the user who created the company
	gorm.Model
}

// ReadCompanyNameByID retrieves the name of the company based on its ID from the database.
func ReadCompanyNameByID(db *gorm.DB, companyID uuid.UUID) (string, error) {
	company := new(Companies)
	check := db.Select("name").Where("id = ?", companyID).First(company)

	if check.Error != nil {
		return "", check.Error
	}

	return company.Name, nil
}
