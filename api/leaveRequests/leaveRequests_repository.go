package leaveRequests

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the LeaveRequests package.
type Database struct {
	DB *gorm.DB
}

// NewLeaveRequestsRepository performs automatic migration of LeaveRequest-related structures in the database.
func NewLeaveRequestsRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.LeaveRequests{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the LeaveRequest structure. Error: ", err)
	}
}

// ReadAll retrieves all LeaveRequests for a specific user based on user ID.
func ReadAll(db *gorm.DB, model domains.LeaveRequests, id uuid.UUID) (domains.LeaveRequests, error) {
	err := db.Where("user_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a LeaveRequests by its unique identifier.
func ReadByID(db *gorm.DB, model domains.LeaveRequests, id uuid.UUID) (domains.LeaveRequests, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of LeaveRequests based on user ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.LeaveRequests, modelID uuid.UUID, limit, offset int) ([]domains.LeaveRequests, error) {
	err := db.Where("user_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}
