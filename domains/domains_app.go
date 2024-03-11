/*

	Package domains provides generic database CRUD (Create, Read, Update, Delete) operations for various data models.

	Structures:
	- UserSessionJWT: Structure for storing user session information in JWT.
		- UserID (uuid.UUID): User ID associated with the session.
		- CompanyID (uuid.UUID): Company ID associated with the session.
		- Roles ([]RolesSessionJWT): List of roles associated with the session.
	- RolesSessionJWT: Structure for storing roles information in JWT.
		- ID (uuid.UUID): Role ID.
		- Name (string): Role name.
		- CompanyID (uuid.UUID): Company ID associated with the role.

	Functions:
	- Create(db *gorm.DB, model any) error: Creates a new record in the database for the provided model.
	- Update(db *gorm.DB, model any, id uuid.UUID) error: Updates a record identified by ID with the provided model in the database.
	- Delete(db *gorm.DB, model any, id uuid.UUID) error: Deletes a record identified by ID for the provided model from the database.
	- CheckByID(db *gorm.DB, model any, id uuid.UUID) error: Checks if a record with the specified ID exists for the provided model in the database.
	- ReadTotalCount(db *gorm.DB, model any, conditionField string, conditionID uuid.UUID) (uint, error): Retrieves the total count of records based on the specified condition.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize generic CRUD operations across different data models.

	Note:
	- The 'model' parameter is expected to be a pointer to the desired data model.
	- The functions return an error if the database operation encounters any issues.

	Last update :
	02/02/2024 12:34

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSessionJWT represents the structure for storing user session information in JWT.
type UserSessionJWT struct {
	CompanyID uuid.UUID `json:"company_id"`
	RoleID    uuid.UUID `json:"role"`
	UserID    uuid.UUID `json:"user_id"`
}

// Create creates a new record in the database for the provided model.
func Create(db *gorm.DB, model interface{}) error {
	return db.Create(model).Error
}

// Update updates a record identified by ID with the provided model in the database.
func Update(db *gorm.DB, model interface{}, id uuid.UUID) error {
	return db.Model(model).Where("id = ?", id).Updates(model).Error
}

// Delete deletes a record identified by ID for the provided model from the database.
func Delete(db *gorm.DB, model interface{}, id uuid.UUID) error {
	return db.Delete(model, id).Error
}

// CheckByID checks if a record with the specified ID exists for the provided model in the database.
func CheckByID(db *gorm.DB, model interface{}, id uuid.UUID) error {
	return db.Select("id").Where("id = ?", id).First(model).Error
}

// ReadTotalCount retrieves the total count of records based on the specified condition.
func ReadTotalCount(db *gorm.DB, model interface{}, conditionField string, conditionID uuid.UUID) (uint, error) {
	var count int64
	check := db.Select("id").Where(conditionField+" = ?", conditionID).Find(model)
	if check.Error != nil {
		return 0, check.Error
	}

	check.Count(&count)
	return uint(count), nil
}

// ReadTotalCount retrieves the total count of records based on the specified condition.
func ReadTotalQuestionsCount(db *gorm.DB, model interface{}, conditionField string, conditionID uuid.UUID) (uint, error) {
	var count int64
	check := db.Select("tests_id").Where(conditionField+" = ?", conditionID).Find(model)
	if check.Error != nil {
		return 0, check.Error
	}

	check.Count(&count)
	return uint(count), nil
}
