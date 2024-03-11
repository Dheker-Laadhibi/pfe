package mission_orders

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
func NewMissionOrdersRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.MissionOrders{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the notification structure. Error: ", err)
	}
}

// ReadAll retrieves all missions orders for a specific user based on user ID.
func ReadAll(db *gorm.DB, model domains.MissionOrders, id uuid.UUID) (domains.MissionOrders, error) {
	err := db.Where("user_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a notification by its unique identifier.
func ReadByID(db *gorm.DB, model domains.MissionOrders, id uuid.UUID) (domains.MissionOrders, error) {
	err := db.First(&model, id).Error
	return model, err
}
// ReadAllPagination retrieves a paginated list of missions based on user ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.MissionOrders, modelID uuid.UUID, limit, offset int) ([]domains.MissionOrders, error) {
	err := db.Where("user_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of missions based on user ID.
func ReadAllList(db *gorm.DB, model []domains.Users, modelID uuid.UUID) ([]domains.Users, error) {
	err := db.Where("user_id = ? ", modelID).Find(&model).Error
	return model, err
}