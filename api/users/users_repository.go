package users

import (
	"labs/domains"
	"math"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the users package.
type Database struct {
	DB *gorm.DB
}

// NewUserRepository performs automatic migration of user-related structures in the database.
func NewUserRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Users{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the user structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of users based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Users, modelID uuid.UUID, limit, offset int) ([]domains.Users, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of users based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Users, modelID uuid.UUID) ([]domains.Users, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a user by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Users, id uuid.UUID) (domains.Users, error) {
	err := db.First(&model, id).Error
	return model, err
}

// GenderPercentages récupère les pourcentages d'hommes et de femmes dans la base de données.
func GenderPercentages(db *gorm.DB, model[]domains.Users ,modelID uuid.UUID) (float64, float64, error) {
    var totalUsers int64
    var maleCount, femaleCount int64

    // Compter le nombre total d'utilisateurs dans la société
    if err := db.Model(&domains.Condidats{}).Where("company_id = ?", modelID).Count(&totalUsers).Error; err != nil {
        return 0, 0, err
    }

    // Compter le nombre d'hommes
    if err := db.Model(&domains.Users{}).Where("company_id = ? AND gender = ?", modelID, "male").Count(&maleCount).Error; err != nil {
        return 0, 0, err
    }

    femaleCount=totalUsers-maleCount;


    // Calculer les pourcentages
    malePercentage := (float64(maleCount) / float64(totalUsers)) * 100
    femalePercentage := (float64(femaleCount) / float64(totalUsers)) * 100

    // Arrondir les pourcentages à l'entier le plus proche
    malePercentage = math.Round(malePercentage)
    femalePercentage = math.Round(femalePercentage)
    
    return malePercentage , femalePercentage, nil

}

