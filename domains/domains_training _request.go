/*

The code defines a struct named TrainingRequest, which represents information about a training request made by an employee in a system. Here's a breakdown of its fields:

ID: This field is of type uuid.UUID, which presumably represents a unique identifier for the training request. It is tagged with gorm metadata specifying that it is the primary key, of type UUID, and cannot be null.

TrainingTitle: This is a string field representing the title of the training request. It cannot be null.

Description: This is a string field representing the description of the training request. It cannot be null.

Reason: This is a string field representing the reason for the training request. It cannot be null.

RequestDate: This is a field of type time.Time representing the date of the training request. It cannot be null.

DecisionCompany: This is a string field representing the decision made by the company regarding the training request. It can be null.

UserID: This is a field of type uuid.UUID representing the user ID associated with the training request.

CompanyID: This is a field of type uuid.UUID representing the company ID associated with the training request.

Each field is tagged with gorm metadata specifying the column name in the database and any additional constraints or properties for the database schema.

*/

package domains

import (
	"time"

	"github.com/google/uuid"
)

// TrainingRequest represents information about training request demanded by an employee  in the system.
type TrainingRequest struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the training Request
	TrainingTitle   string    `gorm:"column:training_title; not null"`             // Object   of the training Request
	Description     string    `gorm:"column:description; not null"`                // Description of   the training Request
	Reason          string    `gorm:"column:reason; not null;"`                    // why to apply on
	RequestDate     time.Time `gorm:"column:request_date; not null;"`
	DecisionCompany string      `gorm:"column:decision_company;"`
	UserID          uuid.UUID `gorm:"column:user_id;"` // User ID associated with the training Request
	CompanyID       uuid.UUID `gorm:"column:company_id;"` //company id

}
