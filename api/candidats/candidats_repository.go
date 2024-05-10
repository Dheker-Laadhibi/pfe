package candidats

import (
	"labs/domains"
	"math"

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

// GenderPercentages récupère les pourcentages d'hommes et de femmes dans la base de données.
func AcceptancePercentages(db *gorm.DB, model[]domains.Condidats,modelID uuid.UUID) (float64, float64, error) {
    var totalCandidates int64
    var accepted, refused int64

    // Compter le nombre total d'utilisateurs dans la société
    if err := db.Model(&domains.Condidats{}).Where("company_id = ?", modelID).Count(&totalCandidates).Error; err != nil {
        return 0, 0, err
    }

    // Compter le nombre d'hommes
    if err := db.Model(&domains.Condidats{}).Where("company_id = ? AND status = ?", modelID, "true").Count(&accepted).Error; err != nil {
        return 0, 0, err
    }

    // Compter le nombre de femmes
    if err := db.Model(&domains.Condidats{}).Where("company_id = ? AND status = ?", modelID, "false").Count(&refused).Error; err != nil {
        return 0, 0, err
    }


    // Calculer les pourcentages
    acceptedPercentage := (float64(accepted) / float64(totalCandidates)) * 100
    refusedPercentage := (float64(refused) / float64(totalCandidates)) * 100

    // Arrondir les pourcentages à l'entier le plus proche
    acceptedPercentage = math.Round(acceptedPercentage)
    refusedPercentage = math.Round(refusedPercentage)
    
    return acceptedPercentage, refusedPercentage, nil
}


// GenderPercentages récupère les pourcentages d'hommes et de femmes dans la base de données.
func levePercentages(db *gorm.DB, model[]domains.Condidats ,modelID uuid.UUID) (float64,float64 ,float64, error) {
    var totalCandidate int64
    var BachelorCount,  MasterCount  int64  
    var otherCount int64

    // Compter le nombre total d'utilisateurs dans la société
    if err := db.Model(&domains.Condidats{}).Where("company_id = ?", modelID).Count(&totalCandidate).Error; err != nil {
        return 0, 0, 0,err
    }

    // Compter le nombre bachelor
    if err := db.Model(&domains.Condidats{}).Where("company_id = ? AND education_level LIKE ?", modelID, "%Bachelor%").Count(&BachelorCount).Error; err != nil {
        return 0, 0, 0 , err
    }


    // Compter le nombre de master
    if err := db.Model(&domains.Condidats{}).Where("company_id = ? AND education_level  LIKE ?", modelID, "%Master%").Count(&MasterCount).Error; err != nil {
        return 0, 0, 0,err
    }


   // Calculer le nombre de candidats avec un niveau d'éducation autre que "Bachelor" ou "Master"
   otherCount = totalCandidate - BachelorCount - MasterCount


    // Calculer les pourcentages
    Bachelor := (float64(BachelorCount) / float64(totalCandidate)) * 100
    Master := (float64(MasterCount) / float64(totalCandidate)) * 100
    Other := (float64(otherCount) / float64(totalCandidate)) * 100
    // Arrondir les pourcentages à l'entier le plus proche
    Bachelor = math.Round(Bachelor)
    Master = math.Round(Master)
    Other= math.Round(Other)
    
    return Bachelor, Master,Other, nil
	
}

