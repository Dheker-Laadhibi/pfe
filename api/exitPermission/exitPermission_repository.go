package exitPermission

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the exitPermission package.
type Database struct {
	DB *gorm.DB
}

// NewExitPermissionRepository performs automatic migration of exitPermission-related structures in the database.
func NewExitPermissionRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.ExitPermission{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the exitPermission structure. Error: ", err)
	}
}

// ReadAll retrieves all exitPermission for a specific user based on user ID.
func ReadAll(db *gorm.DB, model domains.ExitPermission, id uuid.UUID) (domains.ExitPermission, error) {
	err := db.Where("company_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a exitPermission by its unique identifier.
func ReadByID(db *gorm.DB, model domains.ExitPermission, id uuid.UUID) (domains.ExitPermission, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of ExitPermission based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.ExitPermission, modelID uuid.UUID, limit, offset int) ([]domains.ExitPermission, error) {
	err := db.Where("user_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}
