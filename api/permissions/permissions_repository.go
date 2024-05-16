package permissions

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)


/**

IMPORTANT:
The user ID represents the unique identifier of an employee who holds the role of former intern groups.
Please ensure that appropriate permissions and access controls are in place for this user.
*/
// Database represents the database instance for the interns package.
type Database struct {
	DB *gorm.DB
}


// NewInternRepository performs automatic migration of intern-related structures in the database.
func NewPermissionRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Permissions{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the permissions structure. Error: ", err)
	}
}


// ReadAllPagination retrieves a paginated list of interns based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Permissions, modelID uuid.UUID, limit, offset int) ([]domains.Permissions, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}






// ReadAllList retrieves a list of interns based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Permissions, modelID uuid.UUID) ([]domains.Permissions, error) {
	err := db.Where("company_id = ?", modelID).Find(&model).Error
	return model, err
}





// ReadByID retrieves a permission by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Permissions, id uuid.UUID) (domains.Permissions, error) {
	err := db.First(&model, id).Error
	return model, err
}



