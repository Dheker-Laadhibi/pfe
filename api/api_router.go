package api

import (
	"labs/api/auth"
	"labs/api/companies"
	"labs/api/condidats"
	"labs/api/interns"
	"labs/api/mission_orders"
	"labs/api/notifications"
	"labs/api/presences"
	"labs/api/roles"

	"labs/api/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoutesApiInit initializes the API routes for various modules.
func RoutesApiInit(router *gin.Engine, db *gorm.DB) {

	api := router.Group("/api")
	{
		// Initialize authentication routes
		auth.AuthRouterInit(api, db)

		// Initialize user routes
		users.UserRouterInit(api, db)

		// Initialize company routes
		companies.CompanyRouterInit(api, db)

		// Initialize role routes
		roles.RoleRouterInit(api, db)

		// Initialize notification routes
		notifications.NotificationRouterInit(api, db)

		// Initialize Interns routes

		interns.InternRouterInit(api, db)

		// Initialize presences routes
		presences.PresenceRouterInit(api, db)
		// Initialize MissionsOrders routes
		mission_orders.MissionOrdersRouterInit(api, db)

		//  Initialize condidats routes
		condidats.CondidatRouterInit(api, db)

	}
}
