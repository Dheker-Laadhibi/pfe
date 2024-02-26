package notifications

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the notifications package.
type Database struct {
	DB *gorm.DB
}

// NewNotificationRepository performs automatic migration of notification-related structures in the database.
func NewNotificationRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Notifications{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the notification structure. Error: ", err)
	}
}

// ReadAll retrieves all notifications for a specific user based on user ID.
func ReadAll(db *gorm.DB, model domains.Notifications, id uuid.UUID) (domains.Notifications, error) {
	err := db.Where("user_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a notification by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Notifications, id uuid.UUID) (domains.Notifications, error) {
	err := db.First(&model, id).Error
	return model, err
}
