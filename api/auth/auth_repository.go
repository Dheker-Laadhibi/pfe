package auth

import (
	"errors"
	"labs/domains"

	"gorm.io/gorm"
)

// Database represents the database instance for the auth package.
type Database struct {
	DB *gorm.DB
}

// ReadByEmailActive retrieves an active user by email from the database.
func ReadByEmailActive(db *gorm.DB, email string) (*domains.Users, error) {
	var user domains.Users
	err := db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// CheckEmailDuplication checks for email duplication in the database.
func CheckEmailDuplication(db *gorm.DB, email string) error {
	if find := db.Where("email = ?", email).Find(&domains.Users{}); find.Error != nil || find.RowsAffected > 0 {
		return errors.New("email duplication error")
	}
	return nil
}
