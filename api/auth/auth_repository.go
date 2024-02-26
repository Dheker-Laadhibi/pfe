package auth

import (
	"errors"
	"labs/domains"

	"gorm.io/gorm"
)

// Database represents the database instance for the auth package.
type Database struct {
	DB *gorm.DB
}

// ReadByEmailActive retrieves an active user by email from the database.
func ReadByEmailActive(db *gorm.DB, email string) (*domains.Users, error) {
	var user domains.Users
	//, email est une variable contenant l'adresse e-mail de l'utilisateur que l'on souhaite rechercher dans la base de données
	// email  ?:  la recherche doit être effectuée sur la colonne email de la table
	// first récupère le premier enregistrement correspondant à la requête et le stocke dans la variable user.
	err := db.Where("email = ? AND status = true", email).First(&user).Error
	return &user, err
}

// CheckEmailDuplication checks for email duplication in the database.
func CheckEmailDuplication(db *gorm.DB, email string) error {
/*Find pour exécuter la requête de recherche. 
Le deuxième argument de Find est un pointeur vers une instance vide de la structure domains.
Users, indiquant à GORM quel type d'objet il doit rechercher dans la base de données*/
	if find := db.Where("email = ?", email).Find(&domains.Users{}); find.Error != nil || find.RowsAffected > 0 {
		return errors.New("email duplication error")
	}
	// no duplication 
	return nil
}
