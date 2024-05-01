/*

	This Go code snippet defines a package named domains. It imports three packages: "time" for time-related operations, "github.com/google/uuid" for working with UUIDs, and "gorm.io/gorm" which is an ORM (Object-Relational Mapping) library for Go, used for interacting with databases.

Within the package, there is a struct named MissionOrders, which represents information about mission orders in the system. Here's a breakdown of its fields:

ID: This field is of type uuid.UUID, serving as a unique identifier for the mission order. It's tagged with gorm metadata specifying that it's the primary key, of type UUID, and cannot be null.

Object: This is a string field representing the object of the mission order. It cannot be null.

Description: This is a string field representing the description of the mission order. It cannot be null.

StartDate: This is a field of type time.Time representing the start date of the mission order. It cannot be null.

EndDate: This is a field of type time.Time representing the end date of the mission order. It cannot be null.

AdressClient: This is a string field representing the address of the client associated with the mission order. It cannot be null.

Transport: This is a string field representing the transport method for the mission order. It cannot be null.

UserID: This is a field of type uuid.UUID representing the user ID associated with the mission order.

CompanyID: This is a field of type uuid.UUID representing the company ID associated with the mission order.

gorm.Model: This is an embedded struct provided by the GORM library, which includes fields like ID, CreatedAt, UpdatedAt, and DeletedAt to track the model's lifecycle in the database.

Each field is tagged with gorm metadata specifying the column name in the database and any additional constraints or properties for the database schema.







*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MissionOrders represents information about mission_orders  in the system.
type MissionOrders struct {
	ID           uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the missionOrders
	Object       string    `gorm:"column:object; not null"`                  // Object   of the missionOrders
	Description  string    `gorm:"column:description; not null"`                // Description of   the missionOrders
	StartDate    time.Time `gorm:"column:start_date; not null;"`                // StartDate of the missionOrders
	EndDate      time.Time `gorm:"column:end_date; not null;"`                  // EndDte    of the missionOrders
	AdressClient string    `gorm:"column:adress_client; not null;"`             // AdressClient  of the missionOrders
	Transport    string    `gorm:"column:transport; not null;"`                 // Transport of the missionOrders
	UserID       uuid.UUID `gorm:"column:user_id;"`                             // User ID associated with the missionOrders
	CompanyID      uuid.UUID `gorm:"column:company_id;"`                             // company id associated with the missionOrders
	gorm.Model
}



