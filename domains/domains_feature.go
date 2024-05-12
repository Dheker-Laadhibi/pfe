package domains

import (


	"github.com/google/uuid"
	
)

type Feature struct {

	ID            uuid.UUID  `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the exitPermission
	Featurename   string `gorm:"column:feature_name; not null"` 
	CompanyID      uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`       // ID of the company to which the intern belo

}