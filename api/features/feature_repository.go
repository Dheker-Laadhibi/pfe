package features

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



// NewFeatureRepository performs automatic migration of feature-related structures in the database.
func NewFeatureRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Feature{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the Feature structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of roles based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Feature, modelID uuid.UUID, limit, offset int) ([]domains.Feature, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of roles based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Feature, modelID uuid.UUID) ([]domains.Feature, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a role by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Feature, id uuid.UUID) (domains.Feature, error) {
	err := db.First(&model, id).Error
	return model, err
	
}
