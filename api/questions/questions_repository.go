package questions

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the questions package.
type Database struct {
	DB *gorm.DB
}

// NewRoleRepository performs automatic migration of question-related structures in the database.
func NewQuestionRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Questions{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the question structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of questions based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Questions, modelID uuid.UUID, limit, offset int) ([]domains.Questions, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of questions based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Questions, modelID uuid.UUID) ([]domains.Questions, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a question by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Questions, id uuid.UUID) (domains.Questions, error) {
	err := db.First(&model, id).Error
	return model, err
}
