package presences

import (
	"time"

	"github.com/google/uuid"
)

// @Description	PresencesIn represents the input structure for updating the "seen" status of presences.
type PresencesIn struct {
	Matricule uint      `json:"matricule" binding:"required"` // Seen is a boolean indicating whether the presence has been seen or not. It is required.
	Check     time.Time `json:"check" binding:"required"`     // Check is the time of presences
	UserID    uuid.UUID  `json:"userID" binding:"required"`               // unique User ID
} //@name PresencesIn

// @Description	PresencesCount represents the count of presences.
type PresencesCount struct {
	Count uint `json:"count"` // Count is the number of presences.
} //@name PresencesCount

// @Description	PresencesDetails represents detailed information about a specific presence.
type PresencesDetails struct {
	ID        uuid.UUID `json:"id"`              // ID is the unique identifier for the presence.
	Matricule uint      `json:"matricule"`       // Type is the type or category of the presence.
	Check     time.Time `json:"check"`           // Check is the time of presences
	UserID    uuid.UUID  `json:"userID" binding:"required"`  // unique User ID
} //@name PresencesDetails
