package training_request

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the missionsorders package.
type Database struct {
	DB *gorm.DB
}

// NeRepository performs automatic migration of notification-related structures in the database.
func NewNewTrainingRequestRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.TrainingRequest{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the training structure. Error: ", err)
	}
}

// ReadAll retrieves all trainings  for a specific company based oncompanyID.
func ReadAll(db *gorm.DB, model domains.TrainingRequest, id uuid.UUID) (domains.TrainingRequest, error) {
	err := db.Where("user_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a notification by its unique identifier.
func ReadByID(db *gorm.DB, model domains.TrainingRequest, id uuid.UUID) (domains.TrainingRequest, error) {
	err := db.First(&model, id).Error
	return model, err
}
// ReadAllPagination retrieves a paginated list of traiining based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.TrainingRequest, modelID uuid.UUID, limit, offset int) ([]domains.TrainingRequest, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of trainings requests  based on company ID.
func ReadAllList(db *gorm.DB, model []domains.TrainingRequest, modelID uuid.UUID) ([]domains.TrainingRequest, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}



// ReadByID retrieves a user by their unique identifier.
func ReadUserByID(db *gorm.DB, model domains.Users, id uuid.UUID) (domains.Users, error) {
	err := db.First(&model, id).Error
	return model, err
}

