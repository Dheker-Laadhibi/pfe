package companies

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the companies package.
type Database struct {
	DB *gorm.DB
}

// NewCompanyRepository performs automatic migration of company-related structures in the database.
func NewCompanyRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Companies{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the company structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of companies based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Companies, modelID uuid.UUID, limit, offset int) ([]domains.Companies, error) {
	err := db.Where("id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a company by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Companies, id uuid.UUID) (domains.Companies, error) {
	err := db.First(&model, id).Error
	return model, err
}
