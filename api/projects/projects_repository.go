package projects

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the project package.
type Database struct {
	DB *gorm.DB
}

// NewUserRepository performs automatic migration of project-related structures in the database.
func NewProjectRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Project{}, &domains.ProjectsCondidats{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the project structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of projects based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Project, modelID uuid.UUID, limit, offset int) ([]domains.Project, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of projects based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Project, modelID uuid.UUID) ([]domains.Project, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a project by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Project, id uuid.UUID) (domains.Project, error) {
	err := db.First(&model, id).Error
	return model, err
}
