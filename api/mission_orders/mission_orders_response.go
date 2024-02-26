package mission_orders

import (
	"time"

	"github.com/google/uuid"
)

// @Description	MissionOrdersIn represents the input structure for updating the attributes  of MissionOrders.
type MissionOrdersIn struct {
	Object string `json:"object" binding:"required"` // Object   of the missionOrders

	Description string `json:"description" binding:"required"` // Description of   the missionOrders

	AdressClient string `json:"Adress_client" binding:"required"` // Adress client for the Missi

	UserID uuid.UUID `json:"userID" binding:"required"` // ID is the unique identifier for the Missi

} //@name MissionOrdersIn

// @Description	MissionOrdersCount represents the count of MissionOrders.
type MissionOrdersCount struct {
	Count uint `json:"count"` // Count is the number of MissionOrders.
} //@name MissionOrdersCount

// @Description	MissionOrdersDetails represents detailed information about a specific MissionOrders.
type MissionOrdersDetails struct {
	ID          uuid.UUID `json:"id"`          // ID is the unique identifier for the MissionOrders.
	Object      string    `json:"object"`      // Object   of the missionOrders
	Description string    `json:"description"` // Description of   the missionOrders
	Transport   string    `json:"transport"`   // Transport of the missionOrders
	EndDate     time.Time `json:"end_date"`    // EndDte    of the missionOrders
	StartDate   time.Time `json:"start_date"`  // StartDate of the missionOrders
	UserID      uuid.UUID `json:"userID" `     // User ID associated with the missionOrders
} //@name MissionOrdersDetails

// @Description	MissionsOrdersPagination represents the paginated list of missions.
type MissionsPagination struct {
	Items      []MissionsTable `json:"items"`      // Items is a slice containing individual missions details.
	Page       uint            `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint            `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint            `json:"totalCount"` // TotalCount is the total number of missions orders in the entire list.
} //@name MissionsPagination

// @Description	MissionsTable represents a single Missions entry in a table.
type MissionsTable struct {
	ID      uuid.UUID `json:"id"`         // ID is the unique identifier for the missions.
	Object string    `json:"string"` // object is the object  missions of the user.
	Description string    `json:"desciption"` // description is the missions of the user.
	Transport   string    `json:"Transport"`  // Transport of the missions.
	StartDate   time.Time `json:"StartDate"`  // StartDate is the email address of the user.
	EndDate     time.Time `json:"EndDate"`    // EndDate is the timestamp indicating when the missions ends .
} //@name MissionsTable

// @Description	UsersList represents a simplified version of the user for listing purposes.
type MissionsList struct {
	ID     uuid.UUID `json:"id"`     // ID is the unique identifier for the user.
	Object string    `json:"object"` // object is the object of the missions.
} //@name UsersList
