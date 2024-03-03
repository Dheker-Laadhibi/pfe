/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Users: Represents the user information in the system.
		- ID (uuid.UUID): Unique identifier for the user.
		- Firstname (string): The user's first name.
		- Lastname (string): The user's last name.
		- Email (string): User's email address (unique).
		- Password (string): User password.
		- ProfilePicture (string): URL or path to the user's profile picture.
		- Country (string): User's country.
		- Status (bool): User's account status (true for active, false for non-active).
		- LastLogin (time.Time): The last time the user authenticated.
		- Role ([]Roles): The roles assigned to the user.
		- CompanyID (uuid.UUID): ID of the company to which the user belongs.
		- CreatedByUserID (uuid.UUID): ID of the user who created this user.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	- UsersRoles: Represents the roles assigned to users.
		- UserID (uuid.UUID): User's ID.
		- RoleID (uuid.UUID): Role's ID.
		- CompanyID (uuid.UUID): ID of the company associated with the user and role.

	Functions:
	- ReadUsersRoles(db *gorm.DB, userID, companyID uuid.UUID) ([]UsersRoles, error): Reads the roles assigned to a user.
	- CheckEmployeeBelonging(db *gorm.DB, pathCompanyID, sessionUserID, sessionCompanyID uuid.UUID) error: Checks if the user belongs to the specified company.
	- CheckEmployeeSession(db *gorm.DB, pathUserID, sessionUserID, sessionCompanyID uuid.UUID) error: Checks if the user's session matches the specified user and company.

	Dependencies:
	- "errors": Standard Go package for errors handling.
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.
	- "time": Standard Go package for handling time.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling user information in the project.

	Note:
	- The Users structure represents the user information in the system.
	- ReadUsersRoles reads the roles assigned to a user.
	- CheckEmployeeBelonging checks if the user belongs to the specified company.
	- CheckEmployeeSession checks if the user's session matches the specified user and company.

	Last update :
	01/02/2024 10:22

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
	Description  string         `gorm:"column:description; not null;"`               // projects description (mission  )
	Technologies pq.StringArray `gorm:"column:technologies;type:TEXT[]"` // project's technologies
	ExpDate      time.Time      `gorm:"column:exp_date;"`                            // exp date of project
	Condidats    []Condidats     `gorm:"many2many:projects_condidats"`           
	CompanyID    uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the user belongs
	gorm.Model
}

// UsersRoles represents the roles assigned to users.
type ProjectsCondidats struct {
	ProjectID  uuid.UUID `gorm:"column:project_id; primaryKey; type:uuid"`
	CondidatID uuid.UUID `gorm:"column:condidat_id; primaryKey; type:uuid"`
	CompanyID  uuid.UUID `gorm:"column:company_id; type:uuid"`
}

// ReadProjectsCondidats reads the condidats assigned to a project.
func ReadProjectsCondidats(db *gorm.DB, projectID, companyID uuid.UUID) ([]ProjectsCondidats, error) {
	var project []ProjectsCondidats
	err := db.Where("project_id = ? AND company_id = ?", projectID, companyID).Find(&project).Error
	return project, err
}
/* 
// ReadUsersRoles reads the roles assigned to a user.
func ReadUsersRoles(db *gorm.DB, userID, companyID uuid.UUID) ([]UsersRoles, error) {
	var user []UsersRoles
	err := db.Where("users_id = ? AND company_id = ?", userID, companyID).Find(&user).Error
	return user, err
} */


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
