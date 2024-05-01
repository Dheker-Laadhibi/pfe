package presences

import (
	"time"

	"github.com/google/uuid"
)

// @Description	PresencesIn represents the input structure for updating the "seen" status of presences.
type PresencesIn struct {
	Matricule uint      `json:"matricule" binding:"required"` // Seen is a boolean indicating whether the presence has been seen or not. It is required.
	Check     time.Time `json:"check" binding:"required"`     // Check is the time of presences

} //@name PresencesIn

// @Description	PresencesCount represents the count of presences.
type PresencesCount struct {
	Count uint `json:"count"` // Count is the number of presences.
} //@name PresencesCount

// @Description	PresencesDetails represents detailed information about a specific presence.
type PresencesDetails struct {
	ID        uuid.UUID `json:"id"`                        // ID is the unique identifier for the presence.
	Matricule uint      `json:"matricule"`                 // Type is the type or category of the presence.
	Check     time.Time `json:"check"`                     // Check is the time of presences
	UserID    uuid.UUID `json:"userID" binding:"required"` // unique User ID
} //@name PresencesDetails

// @Description	MissionsOrdersPagination represents the paginated list of missions.
type PresencesPagination struct {
	Items      []PresencesTable `json:"items"`      // Items is a slice containing individual missions details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of missions orders in the entire list.
} //@name MissionsPagination

// @Description	presences represents a single presence entry in a table.
type PresencesTable struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the presence.
	Matricule uint      `json:"matricule"` //  matricule is the matricule of the user.
	Check     time.Time `json:"check"`     // presence check.

} //@name MissionsTable

// @Description	UsersList represents a simplified version of the user for listing purposes.
type PresencesList struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the presence.
	Matricule string    `json:"matricule"` // matricule is the matricule of the user.
} //@name UsersList
// @Description	MissionsOrdersPagination represents the paginated list of missions.
