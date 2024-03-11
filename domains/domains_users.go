/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Users: Represents the user information in the system.
		- ID (uuid.UUID): Unique identifier for the user.
		- Firstname (string): The user's first name.
		- Lastname (string): The user's last name.
		- Email (string): User's email address (unique).
		- Password (string): User password.
		- DateOfBirth (time.Time): User's date of birth
		- Gender (string):User's gender
		- Address (string): User's address
		- Country (string): User's country.
		- PhoneNumber (string): User's phone number
		- DateOfHire (time.Time): Date when the user was hired
		- LeaveBalance (int): Balance of leaves for the user
		- CVPath (string): Path to the user's CV
		- LastLogin (time.Time): The last time the user authenticated.
		- DepartementName (string): The name of the department to which the user belongs
		- RoleID (uuid.UUID): The ID of the role assigned to the user
		- CompanyID (uuid.UUID): ID of the company to which the user belongs.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).


	Functions:
	- CheckEmployeeBelonging(db *gorm.DB, pathCompanyID, sessionUserID, sessionCompanyID uuid.UUID) error: Checks if the user belongs to the specified company.
	- CheckEmployeeSession(db *gorm.DB, pathUserID, sessionUserID, sessionCompanyID uuid.UUID) error: Checks if the user's session matches the specified user and company.

	Dependencies:
	- "errors": Standard Go package for errors handling.
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.
	- "time": Standard Go package for handling time.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling user information in the project.

	Note:
	- The Users structure represents the user information in the system.
	- CheckEmployeeBelonging checks if the user belongs to the specified company.
	- CheckEmployeeSession checks if the user's session matches the specified user and company.
	- GetRoleIDByName retrieves the role's ID based on the role name and sessionCompanyID.

	Last update :
	24/02/2024 10:22

*/

package domains

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Users represents the user information in the system.
type Users struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the user
	Firstname       string    `gorm:"column:first_name; not null;"`                // The user's first name
	Lastname        string    `gorm:"column:last_name; not null;"`                 // The user's last name
	Email           string    `gorm:"column:email; not null; unique"`              // User's email address (unique)
	Password        string    `gorm:"column:password; not null;"`                  // User's Password
	DateOfBirth     time.Time `gorm:"column:date_of_birth; not null;"`             // User's date of birth
	Gender          string    `gorm:"column:gender; not null;"`                    // User's gender
	Address         string    `gorm:"column:address; not null;"`                   // User's address
	Country         string    `gorm:"column:country; not null;"`                   // User's country
	PhoneNumber     string    `gorm:"column:phone_number; not null;"`              // User's phone number
	DateOfHire      time.Time `gorm:"column:date_of_hire; not null;"`              // User's date of hire
	LeaveBalance    int       `gorm:"column:leave_balance; not null;"`             // User's leave balance
	CVPath          int       `gorm:"column:cv_path; not null;"`                   // Path to user's CV
	LastLogin       time.Time `gorm:"column:last_login;"`                          // The last time the user authenticated
	DepartementName string    `gorm:"column:departement_name; not null;"`          // User's department name
	RoleID          uuid.UUID `gorm:"column:role_id; type:uuid; not null;"`        // ID of the role assigned to the user
	CompanyID       uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the user belongs
	gorm.Model
}

// CheckEmployeeBelonging checks if the user belongs to the specified company.
func CheckEmployeeBelonging(db *gorm.DB, pathCompanyID, sessionUserID, sessionCompanyID uuid.UUID) error {
	if pathCompanyID != sessionCompanyID {
		return errors.New("error occurred when attempting to verify employee belonging")
	}
	return db.Select("id, company_id").Where("id = ? AND company_id = ?", sessionUserID, pathCompanyID).First(&Users{}).Error
}

// CheckEmployeeSession checks if the user's session matches the specified user and company.
func CheckEmployeeSession(db *gorm.DB, pathUserID, sessionUserID, sessionCompanyID uuid.UUID) error {
	if pathUserID != sessionUserID {
		return errors.New("error occurred when attempting to verify user belonging")
	}
	return db.Select("id, company_id").Where("id = ? AND company_id = ?", pathUserID, sessionCompanyID).First(&Users{}).Error
}

// GetRoleIDByName retrieves the role's ID based on the role name and sessionCompanyID.
func GetRoleIDByName(db *gorm.DB, roleName string, sessionCompanyID uuid.UUID) (uuid.UUID, error) {
	var role Roles

	// Select the role's ID based on the role name and sessionCompanyID
	result := db.Select("id").Where("name = ? AND company_id = ?", roleName, sessionCompanyID).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return uuid.Nil, errors.New("role not found")
		}
		return uuid.Nil, result.Error
	}

	return role.ID, nil
}
