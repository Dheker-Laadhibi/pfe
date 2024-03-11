package interns

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

/**

IMPORTANT:
The user ID represents the unique identifier of an employee who holds the role of former intern groups.
Please ensure that appropriate permissions and access controls are in place for this user.
*/
// Database represents the database instance for the interns package.
type Database struct {
	DB *gorm.DB
}

// NewInternRepository performs automatic migration of intern-related structures in the database.
func NewInternRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Interns{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the intern structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of interns based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Interns, modelID uuid.UUID, limit, offset int) ([]domains.Interns, error) {
	err := db.Where("user_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of interns based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Interns, modelID uuid.UUID) ([]domains.Interns, error) {
	err := db.Where("user_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a intern by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Interns, id uuid.UUID) (domains.Interns, error) {
	err := db.First(&model, id).Error
	return model, err
}
