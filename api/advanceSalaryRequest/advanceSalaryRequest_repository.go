package advanceSalaryRequest

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the advanceSalaryRequest package.
type Database struct {
	DB *gorm.DB
}

// NewAdvanceSalaryRequestRepository performs automatic migration of advanceSalaryRequest-related structures in the database.
func NewAdvanceSalaryRequestsRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.AdvanceSalaryRequests{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the advanceSalaryRequest structure. Error: ", err)
	}
}

// ReadAll retrieves all AdvanceSalaryRequests for a specific user based on user ID.
func ReadAll(db *gorm.DB, model domains.AdvanceSalaryRequests, id uuid.UUID) (domains.AdvanceSalaryRequests, error) {
	err := db.Where("user_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a AdvanceSalaryRequests by its unique identifier.
func ReadByID(db *gorm.DB, model domains.AdvanceSalaryRequests, id uuid.UUID) (domains.AdvanceSalaryRequests, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of AdvanceSalaryRequests based on user ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.AdvanceSalaryRequests, modelID uuid.UUID, limit, offset int) ([]domains.AdvanceSalaryRequests, error) {
	err := db.Where("user_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}
