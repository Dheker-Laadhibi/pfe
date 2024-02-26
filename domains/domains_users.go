/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Users: Represents the user information in the system.
		- ID (uuid.UUID): Unique identifier for the user.
		- Firstname (string): The user's first name.
		- Lastname (string): The user's last name.
		- Email (string): User's email address (unique).
		- Password (string): User password.
		- ProfilePicture (string): URL or path to the user's profile picture.
		- Country (string): User's country.
		- Status (bool): User's account status (true for active, false for non-active).
		- LastLogin (time.Time): The last time the user authenticated.
		- Role ([]Roles): The roles assigned to the user.
		- CompanyID (uuid.UUID): ID of the company to which the user belongs.
		- CreatedByUserID (uuid.UUID): ID of the user who created this user.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	- UsersRoles: Represents the roles assigned to users.
		- UserID (uuid.UUID): User's ID.
		- RoleID (uuid.UUID): Role's ID.
		- CompanyID (uuid.UUID): ID of the company associated with the user and role.

	Functions:
	- ReadUsersRoles(db *gorm.DB, userID, companyID uuid.UUID) ([]UsersRoles, error): Reads the roles assigned to a user.
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
	- ReadUsersRoles reads the roles assigned to a user.
	- CheckEmployeeBelonging checks if the user belongs to the specified company.
	- CheckEmployeeSession checks if the user's session matches the specified user and company.

	Last update :
	01/02/2024 10:22

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
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`     // Unique identifier for the user
	Firstname       string    `gorm:"column:first_name; not null;"`                    // The user's first name
	Lastname        string    `gorm:"column:last_name; not null;"`                     // The user's last name
	Email           string    `gorm:"column:email; not null; unique"`                  // User's email address (unique)
	Password        string    `gorm:"column:password; not null;"`                      // User password
	ProfilePicture  string    `gorm:"column:profile_picture;"`                         // URL or path to the user's profile picture
	Country         string    `gorm:"column:country; not null;"`                       // User's country
	Status          bool      `gorm:"column:status; not null; default:true;"`          // User's account status (true for active, false for non-active)
	LastLogin       time.Time `gorm:"column:last_login;"`                              // The last time the user authenticated
	Role            []Roles   `gorm:"many2many:users_roles"`                           // The roles assigned to the user
	CompanyID       uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`         // ID of the company to which the user belongs
	CreatedByUserID uuid.UUID `gorm:"column:created_by_user_id; type:uuid; not null;"` // ID of the user who created this user
	gorm.Model
}

// UsersRoles represents the roles assigned to users.
type UsersRoles struct {
	UserID    uuid.UUID `gorm:"column:users_id; primaryKey; type:uuid"`
	RoleID    uuid.UUID `gorm:"column:roles_id; primaryKey; type:uuid"`
	CompanyID uuid.UUID `gorm:"column:company_id; type:uuid"`
}

// ReadUsersRoles reads the roles assigned to a user.
func ReadUsersRoles(db *gorm.DB, userID, companyID uuid.UUID) ([]UsersRoles, error) {
	var user []UsersRoles
	err := db.Where("users_id = ? AND company_id = ?", userID, companyID).Find(&user).Error
	return user, err
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
