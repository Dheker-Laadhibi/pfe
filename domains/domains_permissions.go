/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Permissions: Represents the permissions assigned to roles for specific features in the system.
		- ID (uuid.UUID): Unique identifier for the permission.
		- RoleID (uuid.UUID): ID of the role associated with the permission.
		- CompanyID (uuid.UUID): ID of the company associated with the permission.
		- FeatureID (uuid.UUID): ID of the feature for which the permission is granted.
		- FeatureName (string): Name of the feature for better identification.
		- CreatePerm (bool): Indicates whether the role has permission to create for the specified feature.
		- ReadPerm (bool): Indicates whether the role has permission to read for the specified feature.
		- UpdatePerm (bool): Indicates whether the role has permission to update for the specified feature.
		- DeletePerm (bool): Indicates whether the role has permission to delete for the specified feature.
		- CreatedByUserID (uuid.UUID): ID of the user who created the permission.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- CheckPermission(db *gorm.DB, companyID, roleID, featureID uuid.UUID, action string) (bool, error): Checks if the user has permission for a specific action on a resource.

	Dependencies:
	- "errors": Standard Go package for errors.
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling permissions in the project.

	Note:
	- The Permissions structure represents the permissions assigned to roles for specific features in the system.
	- CheckPermission checks if the user has permission for a specific action on a resource.

	Last update :
	01/02/2024 10:22

*/

package domains

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permissions represents the permissions assigned to roles for specific features in the system.
type Permissions struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`     // Unique identifier for the permission
	RoleID          uuid.UUID `gorm:"column:role_id; not null; type:uuid;"`            // ID of the role associated with the permission
	CompanyID       uuid.UUID `gorm:"column:company_id; not null; type:uuid;"`         // ID of the company associated with the permission
	FeatureID       uuid.UUID `gorm:"column:feature_id; not null; type:uuid;"`         // ID of the feature for which the permission is granted
	FeatureName     string    `gorm:"column:feature_name; not null;"`                  // Name of the feature for better identification
	CreatePerm      bool      `gorm:"column:create_perm; not null default:false;"`     // Indicates whether the role has permission to create for the specified feature
	ReadPerm        bool      `gorm:"column:read_perm; not null default:false;"`       // Indicates whether the role has permission to read for the specified feature
	UpdatePerm      bool      `gorm:"column:update_perm; not null default:false;"`     // Indicates whether the role has permission to update for the specified feature
	DeletePerm      bool      `gorm:"column:delete_perm; not null default:false;"`     // Indicates whether the role has permission to delete for the specified feature
	CreatedByUserID uuid.UUID `gorm:"column:created_by_user_id not null;; type:uuid;"` // ID of the user who created the permission
	gorm.Model
}

// CheckPermission checks if the user has permission for a specific action on a resource.
func CheckPermission(db *gorm.DB, companyID, roleID, featureID uuid.UUID, action string) (bool, error) {

	// Init vars
	var permission Permissions

	err := db.Select("role_id, company_id, feature_id, create_perm, read_perm, update_perm, delete_perm").Where("company_id=? AND role_id=? AND feature_id=?", companyID, roleID, featureID).First(&permission).Error
	if err != nil {
		return false, err
	}

	switch action {
	case "create":
		return permission.CreatePerm, nil
	case "read":
		return permission.ReadPerm, nil
	case "update":
		return permission.UpdatePerm, nil
	case "delete":
		return permission.DeletePerm, nil
	default:
		return false, errors.New("invalid action")
	}
}
