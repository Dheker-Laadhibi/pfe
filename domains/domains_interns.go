/*

	Package domains provides the data structures representing entities in the project.
# Interns Struct Documentation

The `Interns` struct represents the user information in the system.

## Fields

- `ID` (`uuid.UUID`): Unique identifier for the user.
- `Firstname` (`string`): The user's first name.
- `Lastname` (`string`): The user's last name.
- `Email` (`string`): User's email address (unique).
- `Password` (`string`): User password.
- `ProfilePicture` (`string`): URL or path to the user's profile picture.
- `Country` (`string`): User's country.
- `Status` (`bool`): User's account status (true for active, false for non-active).
- `LastLogin` (`time.Time`): The last time the user authenticated.
- `UserID` (`uuid.UUID`): ID of the company to which the user belongs.
- `CreatedByUserID` (`uuid.UUID`): ID of the user who created this user.
- `CompanyID` (`uuid.UUID`): ID of the company to which the user belongs.
- `gorm.Model`: GORM's embedded model providing fields like `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`.
 dheker  laadhibiii
## Tags

- `gorm:"column:id; primaryKey; type:uuid; not null;"`: Defines database column properties for the `ID` field.
- `gorm:"column:first_name; not null;"`: Defines database column properties for the `Firstname` field.
- `gorm:"column:last_name; not null;"`: Defines database column properties for the `Lastname` field.
- `gorm:"column:email; not null; unique"`: Defines database column properties for the `Email` field.
- `gorm:"column:password; not null;"`: Defines database column properties for the `Password` field.
- `gorm:"column:profile_picture;"`: Defines database column properties for the `ProfilePicture` field.
- `gorm:"column:country; not null;"`: Defines database column properties for the `Country` field.
- `gorm:"column:status; not null; default:true;"`: Defines database column properties for the `Status` field.
- `gorm:"column:last_login;"`: Defines database column properties for the `LastLogin` field.
- `gorm:"column:user_id; type:uuid; not null;"`: Defines database column properties for the `UserID` field.
- `gorm:"column:created_by_user_id; type:uuid; not null;"`: Defines database column properties for the `CreatedByUserID` field.
- `gorm:"column:company_id; type:uuid; not null;"`: Defines database column properties for the `CompanyID` field.
dheker laadhibiii update

	Last update :
	20/02/2024 16:22


*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/**

IMPORTANT:
The user ID represents the unique identifier of an employee who holds the role of former intern groups.
Please ensure that appropriate permissions and access controls are in place for this user.
*/
//last update 20/02/2024 by dheker

// Interns represents the user information in the system.
type Interns struct {
	ID                         uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the user
	Firstname                  string    `gorm:"column:first_name; not null;"`                // The user's first name
	Lastname                   string    `gorm:"column:last_name; not null;"`                 // The user's last name
	Email                      string    `gorm:"column:email; not null; unique"`              // User's email address (unique)
	Password                   string    `gorm:"column:password; not null;"`
	DateOfBirth                time.Time `gorm:"column:date_of_birth; not null;"` // intern Birthday
	Gender                     string    `gorm:"column:gender; not null;"`        //intern gender
	Adress                     string    `gorm:"column:adress; not null;"`        //intern gender
	PhoneNumber                string    `gorm:"column:phone_number;"`            // intern Phone number
	CountryCode                string    `gorm:"column:country_code; not null;"`  //intern gender
	LevelOfEducation           string    `gorm:"column:education_level; not null;"`
	University                 string    `gorm:"column:university; not null;"` //intern gender
	StartDate                  time.Time    `gorm:"column:start_date; not null;"`
	EndDate                    string    `gorm:"column:end_date; not null;"`
	CvPath                     string    `gorm:"column:cv_path; not null; default:true;"`        // User's account status (true for active, false for non-active)
	LastLogin                  time.Time `gorm:"column:last_login;"`                             // The last time the user authenticated
	EducationalSupervisorName  string    `gorm:"column:educational_supervisor_name;not null;"`  //education supervisor name  from the university
	EducationalSupervisorPhone string    `gorm:"column:educational_supervisor_phone;not null;"` //education supervisor phone from the university
	EducationalSupervisorEmail string    `gorm:"column:educational_supervisor_email;not null;"` //education supervisor email  from the university
	SupervisorID               uuid.UUID `gorm:"column:user_id; type:uuid; not null;"`           // ID of the supervisor to which the intern belongs
	CompanyID                  uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`        // ID of the company to which the intern belongs
	gorm.Model
}

