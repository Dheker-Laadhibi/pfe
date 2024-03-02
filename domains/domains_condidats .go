/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Roles: Represents the roles defined in the system.
		- ID (uuid.UUID): Unique identifier for the role.
		- Name (string): The name of the role.
		- OwningCompanyID (uuid.UUID): ID of the company to which the role belongs.
		- CreatedByUserID (uuid.UUID): ID of the user who created the role.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- ReadRoleName(db *gorm.DB, roleID uuid.UUID) (string, error): Reads the name of the role based on its ID.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling roles in the project.

	Note:
	- The Roles structure represents the roles defined in the system.
	- ReadRoleName reads the name of the role based on its ID.

	Last update :
	01/02/2024 10:22

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Roles represents the roles defined in the system.
type Condidats struct {
	ID             uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the role
	Firstname      string    `gorm:"column:first_name; not null;"`                // The user's first name
	Lastname       string    `gorm:"column:last_name; not null;"`                 // The user's last name
	Email          string    `gorm:"column:email; not null; unique"`              // User's email address (unique)
	Password       string    `gorm:"column:password; not null;"`                  // User password
	University     string    `gorm:"column:university; not null; unique"`         // User's university  address (unique)
	Status         bool      `gorm:"column:status; not null; default:true;"`      // User's account status (true for active, false for non-active)
	Adress         string    `gorm:"column:adress; not null; unique"`             // User's email address (unique)
	Educationlevel string    `gorm:"column:education_level; not null; unique"`    // User's email address (unique)
	CompanyID      uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the user belongs
	gorm.Model
}
