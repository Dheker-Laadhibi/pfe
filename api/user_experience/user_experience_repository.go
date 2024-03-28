package user_experience
import (
	"errors"
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the project package.
type Database struct {
	DB *gorm.DB
}

// NewExperienceRepository performs automatic migration of project-related structures in the database.
func NewExperienceRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.UserExperience{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the  UserExperience structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of users experiences based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.UserExperience, modelID uuid.UUID, limit, offset int) ([]domains.UserExperience, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of projects based on company ID.
func ReadAllList(db *gorm.DB, model []domains.UserExperience, modelID uuid.UUID) ([]domains.UserExperience, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a project by their unique identifier.
func ReadByID(db *gorm.DB, model domains.UserExperience, id uuid.UUID) (domains.UserExperience, error) {
	err := db.First(&model, id).Error
	return model, err
}






// CheckUserExists vérifie si un utilisateur avec un certain ID existe dans la base de données.
func CheckUserExists(db *gorm.DB, userID uuid.UUID) bool {
    var user domains.UserExperience
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        // L'utilisateur n'existe pas dans la base de données
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return true
        }
       
    }
    // L'utilisateur existe dans la base de données
    return false
}
