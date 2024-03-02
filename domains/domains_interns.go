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
	ID               uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the intern
	Firstname        string    `gorm:"column:first_name; not null;"`                // The intern's first name
	Lastname         string    `gorm:"column:last_name; not null;"`                 // The intern's last name
	Email            string    `gorm:"column:email; not null; unique"`              // intern's email address (unique)
	Password         string    `gorm:"column:password; not null;"`
	DateOfBirth      time.Time `gorm:"column:date_of_birth; not null;"`   // intern Birthday
	Gender           string    `gorm:"column:gender; not null;"`          //intern's gender
	Adress           string    `gorm:"column:adress; not null;"`          //intern's adress
	PhoneNumber      string    `gorm:"column:phone_number;"`              // intern's Phone number
	CountryCode      string    `gorm:"column:country_code; not null;"`    //intern's country code
	LevelOfEducation string    `gorm:"column:education_level; not null;"` //intern's educqtion lvl
	University       string    `gorm:"column:university; not null;"`      //intern's unviversity
	StartDate        time.Time `gorm:"column:start_date; not null;"`      // start date of the internship

	EndDate                    string    `gorm:"column:end_date; not null;"`                    // end date  date of the internship
	CvPath                     string    `gorm:"column:cv_path; not null; default:true;"`       // User's account status (true for active, false for non-active)
	LastLogin                  time.Time `gorm:"column:last_login;"`                            // The last time the user authenticated
	EducationalSupervisorName  string    `gorm:"column:educational_supervisor_name;not null;"`  //education supervisor name  from the university
	EducationalSupervisorPhone string    `gorm:"column:educational_supervisor_phone;not null;"` //education supervisor phone from the university
	EducationalSupervisorEmail string    `gorm:"column:educational_supervisor_email;not null;"` //education supervisor email  from the university
	SupervisorID               uuid.UUID `gorm:"column:user_id; type:uuid; not null;"`          // ID of the supervisor to which the intern belongs
	CompanyID                  uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`       // ID of the company to which the intern belongs
	gorm.Model
}
