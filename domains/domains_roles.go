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
type Roles struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`     // Unique identifier for the role
	Name            string    `gorm:"column:name; not null;"`                          // The name of the role
	OwningCompanyID uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`         // ID of the company to which the role belongs
	CreatedByUserID uuid.UUID `gorm:"column:created_by_user_id; type:uuid; not null;"` // ID of the user who created the role
	gorm.Model
}


// ReadRoleName reads the name of the role based on its ID.
func ReadRoleName(db *gorm.DB, roleID uuid.UUID) (string, error) {
	role := new(Roles)
	err := db.Select("id, name").Where("id = ?", roleID).First(role).Error
	return role.Name, err
}







// GetRoleIDByName retrieves the ID of a role by its name.
/*func GetRoleIDByName(db *gorm.DB, roleName string) (uuid.UUID, error) {
      role := new(Roles)
    err := db.Where("name = ?", roleName).First(&role).Error
    if err != nil {
        return uuid.Nil, err
    }
    return role.ID, nil
}*/