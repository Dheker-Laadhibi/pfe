package roles

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the roles package.
type Database struct {
	DB *gorm.DB
}

// NewRoleRepository performs automatic migration of role-related structures in the database.
func NewRoleRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Roles{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the role structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of roles based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Roles, modelID uuid.UUID, limit, offset int) ([]domains.Roles, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of roles based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Roles, modelID uuid.UUID) ([]domains.Roles, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a role by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Roles, id uuid.UUID) (domains.Roles, error) {
	err := db.First(&model, id).Error
	return model, err
}
