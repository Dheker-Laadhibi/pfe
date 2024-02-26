package users

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)
//last update 20/02/2024 by dheker 
// Database represents the database instance for the users package.
type Database struct {
	DB *gorm.DB
}

// NewUserRepository performs automatic migration of user-related structures in the database.
func NewUserRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Users{}, &domains.UsersRoles{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the user structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of users based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Users, modelID uuid.UUID, limit, offset int) ([]domains.Users, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of users based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Users, modelID uuid.UUID) ([]domains.Users, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a user by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Users, id uuid.UUID) (domains.Users, error) {
	err := db.First(&model, id).Error
	return model, err
}
