package candidats

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the roles package.
type Database struct {
	DB *gorm.DB
}

// NewCondidatRepository performs automatic migration of candidat structures in the database.
func NewCondidatRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Condidats{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the candidats structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of candidat based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Condidats, modelID uuid.UUID, limit, offset int) ([]domains.Condidats, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of candidat based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Condidats, modelID uuid.UUID) ([]domains.Condidats, error) {
	err := db.Where("company_id= ? ", modelID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a candidat by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Condidats, id uuid.UUID) (domains.Condidats, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadByEmailActive retrieves an active Candidat by email from the database.
func ReadByEmailActive(db *gorm.DB, email string) (*domains.Condidats, error) {
	var Candidat domains.Condidats
	//, email est une variable contenant l'adresse e-mail de Candidat que l'on souhaite rechercher dans la base de données
	// email  ?:  la recherche doit être effectuée sur la colonne email de la table
	// first récupère le premier enregistrement correspondant à la requête et le stocke dans la variable Candidat.
	err := db.Where("email = ? ", email).First(&Candidat).Error
	return &Candidat, err
}
