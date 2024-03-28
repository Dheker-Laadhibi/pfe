package loanRequests

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the LoanRequests package.
type Database struct {
	DB *gorm.DB
}

// NewLoanRequestsRepository performs automatic migration of LoanRequests-related structures in the database.
func NewLoanRequestRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.LoanRequests{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the LoanRequests structure. Error: ", err)
	}
}

// ReadAll retrieves all LoanRequests for a specific user based on user ID.
func ReadAll(db *gorm.DB, model domains.LoanRequests, id uuid.UUID) (domains.LoanRequests, error) {
	err := db.Where("company_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a LoanRequests by its unique identifier.
func ReadByID(db *gorm.DB, model domains.LoanRequests, id uuid.UUID) (domains.LoanRequests, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of LoanRequests based on user ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.LoanRequests, conditionField string, modelID uuid.UUID, limit, offset int) ([]domains.LoanRequests, error) {
	err := db.Where(conditionField+" = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}
