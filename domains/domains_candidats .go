/*
the struct named Condidats, which seems to represent candidates in a system. Here's a breakdown of its fields:

ID: This field is of type uuid.UUID, presumably serving as a unique identifier for the candidate. It's tagged with gorm metadata specifying that it's the primary key, of type UUID, and cannot be null.

Firstname: This is a string field representing the candidate's first name. It cannot be null.

Lastname: This is a string field representing the candidate's last name. It cannot be null.

Email: This is a string field representing the candidate's email address. It must be unique in the database.

Password: This is a string field representing the candidate's password. It cannot be null.

University: This is a string field representing the candidate's university. It cannot be null.

Status: This is a boolean field representing the account status of the candidate. It's true for active and false for non-active, with a default value of true.

Adress: This is a string field representing the candidate's address. It cannot be null.

Educationlevel: This is a string field representing the candidate's education level. It cannot be null.

RoleID: This is a field of type uuid.UUID representing the ID of the role associated with the candidate.

CompanyID: This is a field of type uuid.UUID representing the ID of the company to which the candidate belongs. It cannot be null.

gorm.Model: This is an embedded struct provided by the GORM library, which includes fields like ID, CreatedAt, UpdatedAt, and DeletedAt to track the model's lifecycle in the database.

Each field is tagged with gorm metadata specifying the column name in the database and any additional constraints or properties for the database schema.







*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Roles represents the roles defined in the system.
type Condidats struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the role
	Firstname string    `gorm:"column:first_name; not null;"`                // The user's first name
	Lastname  string    `gorm:"column:last_name; not null;"`                 // The user's last name
	Email     string    `gorm:"column:email; not null; unique"`              // User's email address (unique)
	Password  string    `gorm:"column:password; not null;"`                  // User password
	University     string    `gorm:"column:university; not null; not null"`      // User's university  address (unique)
	Status         bool      `gorm:"column:status; not null; default:false;"`   // User's account status (true for active, false for non-active)
	Adress         string    `gorm:"column:adress; not null; not null "`        // User's email address (unique)
	Educationlevel string    `gorm:"column:education_level; not null;"` // User's email address (unique)
	RoleID uuid.UUID `gorm:"column:role_id;type:uuid"`
	CompanyID      uuid.UUID `gorm:"column:company_id; type:uuid; not null;"` // ID of the company to which the user belongs
	gorm.Model
}
