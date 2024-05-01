package notifications

import (
	"labs/constants"
	"labs/domains"
	"net/http"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ReadNotifications	Handles the retrieval of all notifications.
// @Summary        		Get notifications
// @Description    		Get all notifications.
// @Tags				Notifications
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID				path			string		true		"User ID"
// @Success				200					{array}			notifications.NotificationsDetails
// @Failure				400					{object}		utils.ApiResponses		"Invalid request"
// @Failure				401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			   	500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/notifications/{userID}	[get]
func (db Database) ReadNotifications(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the notification ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all notification data from the database
	notifications, err := ReadAll(db.DB, domains.Notifications{}, session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, notifications)
}

// ReadNotificationsCount	Handles the retrieval the number of all notifications.
// @Summary        			Get notifications count
// @Description    			Get all notifications count.
// @Tags					Notifications
// @Produce					json
// @Security 				ApiKeyAuth
// @Param					userID				path			string		true		"User ID"
// @Success					200					{object}		notifications.NotificationsCount
// @Failure					400					{object}		utils.ApiResponses		"Invalid request"
// @Failure					401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure					403					{object}		utils.ApiResponses		"Forbidden"
// @Failure					500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router					/notifications/{userID}/count	[get]
func (db Database) ReadNotificationsCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the notification ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all notification data from the database
	notifications, err := domains.ReadTotalCount(db.DB, &domains.Notifications{}, "user_id", session.UserID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	count := NotificationsCount{
		Count: notifications,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, count)
}

// ReadNotification		Handles the retrieval of one notification.
// @Summary        		Get notification
// @Description    		Get one notification.
// @Tags				Notifications
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string		true		"User ID"
// @Param				ID						path			string		true		"Notification ID"
// @Success				200						{object}		notifications.NotificationsDetails
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/notifications/{userID}/{ID}	[get]
func (db Database) ReadNotification(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the notification ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve notification data from the database
	notification, err := ReadByID(db.DB, domains.Notifications{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving notification data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a user structure as a response
	details := NotificationsDetails{
		ID:        notification.ID,
		Type:      notification.Type,
		Content:   notification.Content,
		Seen:      notification.Seen,
		CreatedAt: notification.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateNotification 	Handles the update of a notification.
// @Summary        		Update notification
// @Description    		Update one notification.
// @Tags				Notifications
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string							true		"User ID"
// @Param				ID						path			string							true		"Notification ID"
// @Param				request					body			notifications.NotificationsIn	true		"Notification query params"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses		"Invalid request"
// @Failure				401						{object}		utils.ApiResponses		"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses		"Forbidden"
// @Failure				500						{object}		utils.ApiResponses		"Internal Server Error"
// @Router				/notifications/{userID}/{ID}	[put]
func (db Database) UpdateNotification(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the notification ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a NotificationIn struct
	notification := new(NotificationsIn)
	if err := ctx.ShouldBindJSON(notification); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the notification with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Notifications{}, objectID); err != nil {
		logrus.Error("Error checking if the notification with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the notification data in the database
	dbNotification := &domains.Notifications{
		Seen: notification.Seen,
	}
	if err = domains.Update(db.DB, dbNotification, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteNotification 	Handles the deletion of a notification.
// @Summary        		Delete notification
// @Description    		Delete one notification.
// @Tags				Notifications
// @Accept				json
// @Produce				json
// @Security 			ApiKeyAuth
// @Param				userID					path			string			true		"Company ID"
// @Param				ID						path			string			true		"Notification ID"
// @Success				200						{object}		utils.ApiResponses
// @Failure				400						{object}		utils.ApiResponses			"Invalid request"
// @Failure				401						{object}		utils.ApiResponses			"Unauthorized"
// @Failure				403						{object}		utils.ApiResponses			"Forbidden"
// @Failure				500						{object}		utils.ApiResponses			"Internal Server Error"
// @Router				/notifications/{userID}/{ID}	[delete]
func (db Database) DeleteNotification(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the user ID from the request parameter
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the notification ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeSession(db.DB, userID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Delete the user data from the database
	if err := domains.Delete(db.DB, &domains.Notifications{}, objectID); err != nil {
		logrus.Error("Error deleting user data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
