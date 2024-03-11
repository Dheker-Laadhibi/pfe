package notifications

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NotificationRouterInit initializes the routes related to notifications.
func NotificationRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewNotificationRepository(db)

	// Private
	notifications := router.Group("/notifications/:userID")
	{

		// GET endpoint to retrieve all notifications for a specific user
		notifications.GET("", baseInstance.ReadNotifications)

		// GET endpoint to retrieve the count of notifications for a specific user
		notifications.GET("/count", baseInstance.ReadNotificationsCount)

		// GET endpoint to retrieve details of a specific notification for a specific user
		notifications.GET("/:ID", baseInstance.ReadNotification)

		// PUT endpoint to update the details of a specific notification for a specific user
		notifications.PUT("/:ID", baseInstance.UpdateNotification)

		// DELETE endpoint to delete a specific notification for a specific user
		notifications.DELETE("/:ID", baseInstance.DeleteNotification)
	}
}
