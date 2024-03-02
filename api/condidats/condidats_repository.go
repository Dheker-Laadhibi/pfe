package condidats

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

// NewCondidatRepository performs automatic migration of condidats structures in the database.
func NewCondidatRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Condidats{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the condidats structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of roles based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Condidats, modelID uuid.UUID, limit, offset int) ([]domains.Condidats, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of condidats based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Condidats, modelID uuid.UUID) ([]domains.Condidats, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a condidat by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Condidats, id uuid.UUID) (domains.Condidats, error) {
	err := db.First(&model, id).Error
	return model, err
}
