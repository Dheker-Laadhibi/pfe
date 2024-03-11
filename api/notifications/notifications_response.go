package notifications

import (
	"time"

	"github.com/google/uuid"
)

// @Description	NotificationsIn represents the input structure for updating the "seen" status of notifications.
type NotificationsIn struct {
	Seen bool `json:"seen" binding:"required"` // Seen is a boolean indicating whether the notification has been seen or not. It is required.
} //@name NotificationsIn

// @Description	NotificationsCount represents the count of notifications.
type NotificationsCount struct {
	Count uint `json:"count"` // Count is the number of notifications.
} //@name NotificationsCount

// @Description	NotificationsDetails represents detailed information about a specific notification.
type NotificationsDetails struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier for the notification.
	Type      string    `json:"type"`      // Type is the type or category of the notification.
	Content   string    `json:"content"`   // Content is the textual content of the notification.
	Seen      bool      `json:"seen"`      // Seen is a boolean indicating whether the notification has been seen or not.
	CreatedAt time.Time `json:"createdAt"` // CreatedAt is the timestamp indicating when the notification was created.
} //@name NotificationsDetails
