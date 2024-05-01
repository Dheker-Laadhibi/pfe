package presences

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

// NewPresenceRepository performs automatic migration of notification-related structures in the database.
func NewPresenceRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Presences{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the notification structure. Error: ", err)
	}
}

// ReadAll retrieves all presences for a specific presence  based on user ID.
func ReadAll(db *gorm.DB, model domains.Presences, id uuid.UUID) (domains.Presences, error) {
	err := db.Where("user_id = ?", id).Find(&model).Error
	return model, err
}

// ReadByID retrieves a presence by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Presences, id uuid.UUID) (domains.Presences, error) {
	err := db.First(&model, id).Error
	return model, err
}
// ReadAllPagination retrieves a paginated list of presences based on user ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Presences, modelID uuid.UUID, limit, offset int) ([]domains.Presences, error) {
	err := db.Where("user_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of presences based on user ID.
func ReadAllList(db *gorm.DB, model []domains.Presences, modelID uuid.UUID) ([]domains.Presences, error) {
	err := db.Where("user_id = ? ", modelID).Find(&model).Error
	return model, err
}