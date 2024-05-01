/*


This Go code snippet defines a package named domains. It imports two packages: "github.com/google/uuid" for working with UUIDs and "gorm.io/gorm" which is an ORM (Object-Relational Mapping) library for Go, used for interacting with databases.

Within the package, there is a struct named UserExperience, which represents the professional experiences of a user. Here's a breakdown of its fields:

ID: This field is of type uuid.UUID, serving as a unique identifier for the user experience. It's tagged with gorm metadata specifying that it's the primary key, of type UUID, and cannot be null.

ProfessionalTraining: This is a string field representing the professional training or experience of the user. It cannot be null.

UserID: This is a field of type uuid.UUID representing the ID of the user who owns the experiences.

CompanyID: This is a field of type uuid.UUID representing the ID of the company to which the user belongs. It cannot be null.

gorm.Model: This is an embedded struct provided by the GORM library, which includes fields like ID, CreatedAt, UpdatedAt, and DeletedAt to track the model's lifecycle in the database.

Each field is tagged with gorm metadata specifying the column name in the database and any additional constraints or properties for the database schema.

*/

package domains

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// UserExperience represents the user professional experiences.
type UserExperience struct {
	ID                   uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the   UserExperience
	ProfessionalTraining string    `gorm:"column:professional_trainings;not null;"`     //experience

	UserID    uuid.UUID `gorm:"column:user_id;"`                         //  the use that owns the  experiences
	CompanyID uuid.UUID `gorm:"column:company_id; type:uuid; not null;"` // ID of the company to which the user  belongs
	gorm.Model
}
