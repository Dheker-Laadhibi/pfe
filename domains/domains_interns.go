package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
	Package domains provides the data structures representing entities in the project.
	Structures:
	- Interns: Represents the intern information in the system.
		- ID (uuid.UUID): Unique identifier for the intern.
		- Firstname (string): The intern's first name.
		- Lastname (string): The intern's last name.
		- Email (string): Intern's email address (unique).
		- Password (string): Intern's password.
		- DateOfBirth (time.Time): Intern's date of birth.
		- Gender (string): Intern's gender.
		- Address (string): Intern's address.
		- PhoneNumber (string): Intern's phone number.
		- CountryCode (string): Intern's country code.
		- LevelOfEducation (string): Intern's education level.
		- University (string): Intern's university.
		- StartDate (time.Time): Start date of the internship.
		- EndDate (string): End date of the internship.
		- CvPath (string): Path to intern's curriculum vitae.
		- LastLogin (time.Time): The last time the intern authenticated.
		- EducationalSupervisorName (string): Education supervisor name from the university.
		- EducationalSupervisorPhone (string): Education supervisor phone from the university.
		- EducationalSupervisorEmail (string): Education supervisor email from the university.
		- SupervisorID (uuid.UUID): ID of the supervisor to which the intern belongs.
		- CompanyID (uuid.UUID): ID of the company to which the intern belongs.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Dependencies:
	- "time": Standard Go package for handling time.
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.
	Usage:
	- Import this package to utilize the provided data structures for handling intern information in the project.
	Note:
	- The Interns structure represents the intern information in the system.
	Last update:
	24/02/2024 10:22
*/

type Interns struct {
	ID               uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the intern
	Firstname        string    `gorm:"column:first_name; not null;"`                // The intern's first name
	Lastname         string    `gorm:"column:last_name; not null;"`                 // The intern's last name
	Email            string    `gorm:"column:email; not null; unique"`              // intern's email address (unique)
	Password         string    `gorm:"column:password; not null;"`
	DateOfBirth      time.Time `gorm:"column:date_of_birth; not null;"`   // intern's Birthday
	Gender           string    `gorm:"column:gender; not null;"`          //intern's gender
	Adress           string    `gorm:"column:adress; not null;"`          //intern's adress
	PhoneNumber      string    `gorm:"column:phone_number;"`              // intern's Phone number
	CountryCode      string    `gorm:"column:country_code; not null;"`    //intern's country code
	LevelOfEducation string    `gorm:"column:education_level; not null;"` //intern's educqtion lvl
	University       string    `gorm:"column:university; not null;"`      //intern's unviversity
	StartDate        time.Time `gorm:"column:start_date; not null;"`      // start date of the internship
	EndDate          time.Time    `gorm:"column:end_date; not null;"`                    // end date  date of the internship
	CvPath                     string    `gorm:"column:cv_path; not null; default:true;"`       // intern 's   curriculum vitae
	LastLogin                  time.Time `gorm:"column:last_login;"`                            // The last time the intern authenticated
	EducationalSupervisorName  string    `gorm:"column:educational_supervisor_name;not null;"`  //education supervisor name  from the university
	EducationalSupervisorPhone string    `gorm:"column:educational_supervisor_phone;not null;"` //education supervisor phone from the university
	EducationalSupervisorEmail string    `gorm:"column:educational_supervisor_email;not null;"` //education supervisor email  from the university
	SupervisorID               uuid.UUID `gorm:"column:user_id; type:uuid; not null;"`          // ID of the supervisor to which the intern belongs
	CompanyID                  uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`       // ID of the company to which the intern belongs
	gorm.Model
}
