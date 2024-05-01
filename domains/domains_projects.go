/*

	This package provides data structures representing entities in the project along with functions to manipulate these data. Below is a detailed description of the structures and functions provided:

Data Structures:
Project
ID (uuid.UUID): Unique identifier for the project.
Code (string): The project's code.
Projectname (string): The project's name.
Specialty (string): The project's specialty.
Description (string): Description of the project.
Technologies (pq.StringArray): Technologies used in the project.
ExpDate (time.Time): Expiration date of the project.
Condidats ([]Condidats): Candidates associated with the project.
CompanyID (uuid.UUID): ID of the company to which the project belongs.
gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).
ProjectsCondidats
ProjectID (uuid.UUID): ID of the project.
CondidatID (uuid.UUID): ID of the candidate.
CompanyID (uuid.UUID): ID of the company associated with the project and candidate.
gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).
Functions:
ReadProjectsCondidats(db *gorm.DB, projectID, companyID uuid.UUID) ([]ProjectsCondidats, error): Reads the candidates assigned to a project.
CheckProjectCodeExists(db *gorm.DB, code string) (bool, error): Checks if a project code already exists in the database.
*AssignProjectCondidat(db gorm.DB, projectID, condidatID, companyID uuid.UUID) error: Associates a candidate with a project in the ProjectsCondidats table.
FindProjectIDByCode(db *gorm.DB, projectCode string) (string, error): Searches for the project ID based on the project code in the database.
Dependencies:
"github.com/google/uuid": Package for working with UUIDs.
"github.com/lib/pq": Package for manipulating string arrays in PostgreSQL.
"gorm.io/gorm": The GORM library for object-relational mapping in Go.
"time": Standard Go package for handling time.
Usage:
Import this package to utilize the provided data structures and functions for manipulating project information in the project.
Note:
The Project structure represents project information in the system.
ReadProjectsCondidats reads the candidates assigned to a project.
CheckProjectCodeExists checks if a project code already exists in the database.
AssignProjectCondidat associates a candidate with a project.
FindProjectIDByCode searches for the project ID based on the project code.
Last update:
01/02/2024 10:22 - dheker

*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Users represents the user information in the system.
type Project struct {
	ID           uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the pfe project
	Code         string         `gorm:"column:code; not null;"`                      // The project code
	Projectname  string         `gorm:"column:first_name; not null;"`                // The project name
	Specialty    string         `gorm:"column:specialty; not null;"`                 // The specialty of the project
	Description  string         `gorm:"column:description; not null;"`               // projects description (mission  )
	Technologies pq.StringArray `gorm:"column:technologies;type:TEXT[]"`             // project's technologies
	ExpDate      time.Time      `gorm:"column:exp_date;"`                            // exp date of project
	Condidats    []Condidats    `gorm:"many2many:projects_condidats"`                //projects candidats 
	CompanyID    uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the user belongs
	gorm.Model     
}

// UsersRoles represents the roles assigned to users.
type ProjectsCondidats struct {
	ProjectID  uuid.UUID `gorm:"column:project_id; primaryKey; type:uuid"`
	CondidatID uuid.UUID `gorm:"column:condidat_id; primaryKey; type:uuid"`
	CompanyID  uuid.UUID `gorm:"column:company_id; type:uuid"`
	gorm.Model
}

// ReadProjectsCondidats reads the condidats assigned to a project.
func ReadProjectsCondidats(db *gorm.DB, projectID, companyID uuid.UUID) ([]ProjectsCondidats, error) {
	var project []ProjectsCondidats
	err := db.Where("project_id = ? AND company_id = ?", projectID, companyID).Find(&project).Error
	return project, err
}



// CheckProjectCodeExists vérifie si un code de projet existe déjà dans la base de données.
func CheckProjectCodeExists(db *gorm.DB, code string) (bool, error) {
	// Déclarez une variable pour stocker le nombre de projets avec le code donné.
	var count int64

	// Exécutez une requête pour compter le nombre de projets avec le code donné.
	if err := db.Model(&Project{}).Where("code = ?", code).Count(&count).Error; err != nil {
		// Si une erreur se produit lors de l'exécution de la requête, retournez l'erreur.
		return false, err
	}

	// Si le nombre de projets avec le code donné est supérieur à zéro, cela signifie que le code de projet existe déjà.
	// Sinon, le code de projet n'existe pas encore.
	return count > 0, nil
}

// associe un condidat à un projet dans la table ProjectsCondidats.
func AssignProjectCondidat(db *gorm.DB, projectID, condidatID, companyID uuid.UUID) error {
	// Créez une instance de ProjectsCondidats avec les identifiants fournis.
	projectCondidat := ProjectsCondidats{
		ProjectID:  projectID,
		CondidatID: condidatID,
		CompanyID:  companyID,
	}

	if err := db.Create(&projectCondidat).Error; err != nil {

		return err
	}

	// Retournez nil pour indiquer qu'il n'y a pas eu d'erreur.
	return nil
}

// FindProjectIDByCode recherche l'ID du projet basé sur le code du projet dans la base de données.
func FindProjectIDByCode(db *gorm.DB, projectCode string) (string, error) {
	var projectID string
	// Exécutez une requête pour trouver l'ID du projet basé sur le code du projet donné
	if err := db.Model(&Project{}).Where("code = ?", projectCode).Pluck("id", &projectID).Error; err != nil {
		// Si une erreur se produit lors de l'exécution de la requête, renvoyez l'erreur.
		return "", err
	}
	return projectID, nil
}
