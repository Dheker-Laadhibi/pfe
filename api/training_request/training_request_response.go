package training_request

import (
	"time"

	"github.com/google/uuid"
)

// @Description	TrainingRequestIn represents the input structure for updating the attributes  of MissionOrders.
type TrainingRequestIn struct {
	TrainingTitle string `json:"training_title" binding:"required"` // Object   of the missionOrders

	Description string `json:"description" binding:"required"` // Description of   the missionOrders

	Reason string `json:"reason" binding:"required"` // Adress client for the Missi

	//RequestDate
	RequestDate  string   `json:"request_date" binding:"required"`
	DecisionCompany string  `json:"decision_company" `
	UserID uuid.UUID `json:"userID" binding:"required"` // ID is the unique identifier for the Missi

} //@name TrainingRequestIn

// @Description	MissionOrdersCount represents the count of MissionOrders.
type TrainingRequestsCount struct {
	Count uint `json:"count"` // Count is the number of MissionOrders.
} //@name MissionOrdersCount

// @Description	TrainingRequestDetails represents detailed information about a specific MissionOrders.
type TrainingRequestDetails struct {
	ID           uuid.UUID      `json:"id"` 
	TrainingTitle string `json:"training_title"` // Object   of the missionOrders

	Description string `json:"description" ` // Description of   the missionOrders

	Reason string `json:"reason"` // Adress client for the Missi
	//RequestDate
	RequestDate   time.Time     `json:"request_date"`
	DecisionCompany string  `json:"decision_company"`
	UserID uuid.UUID `json:"userID"` // ID is the unique identifier for the Missi

} //@name TrainingRequestDetails

// @Description	TrainingRequestPagination represents the paginated list of missions.
type TrainingRequestPagination struct {
	Items      []TrainingRequestTable `json:"items"`      // Items is a slice containing individual missions details.
	Page       uint            `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint            `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint            `json:"totalCount"` // TotalCount is the total number of missions orders in the entire list.
} //@name TrainingRequestPagination

// @Description	TrainingRequestTable represents a single TrainingRequest entry in a table.
type TrainingRequestTable struct {
	ID           uuid.UUID      `json:"id"` 
	TrainingTitle string `json:"training_title"` // Object   of the missionOrders

	Description string `json:"description" ` // Description of   the missionOrders

	Reason string `json:"reason"` // Adress client for the Missi
	//RequestDate
	RequestDate   time.Time     `json:"request_date"`
	DecisionCompany string  `json:"decision_company"`
	UserID uuid.UUID `json:"userID" ` // ID is the unique identifier for the Missi

} //@name TrainingRequestTable

// @Description	TrainingRequestList represents a simplified version of the TrainingRequest for listing purposes.
type TrainingRequestList struct {
	ID     uuid.UUID `json:"id"`     // ID is the unique identifier for the user.
	TrainingTitle string    `json:"training_title"` // object is the object of the missions.
} //@name TrainingRequestList
// @Description	TrainingRequestDescision represents the descision of training request .
type TrainingRequestDescision struct {
	DecisionCompany string  `json:"decision_company"`
} //@name TrainingRequestDescision
