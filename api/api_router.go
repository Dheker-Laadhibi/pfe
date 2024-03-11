package api

import (
	"labs/api/advanceSalaryRequest"
	"labs/api/auth"
	"labs/api/candidats"
	"labs/api/companies"
	"labs/api/exitPermission"
	"labs/api/interns"
	"labs/api/leaveRequests"
	"labs/api/loanRequests"
	"labs/api/mission_orders"
	"labs/api/notifications"
	"labs/api/presences"
	"labs/api/projects"
	"labs/api/questions"
	"labs/api/roles"
	"labs/api/tests"
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

		// Initialize leave_requests routes
		leaveRequests.LeaveRouterInit(api, db)

		// Initialize exitPermission routes
		exitPermission.ExitPermissionRouterInit(api, db)

		// Initialize advanceSalaryRequest routes
		advanceSalaryRequest.AdvanceSalaryRequestsRouterInit(api, db)

		// Initialize loanRequests routes
		loanRequests.LoanRequestRouterInit(api, db)

		// Initialize Interns routes
		interns.InternRouterInit(api, db)

		// Initialize presences routes
		presences.PresenceRouterInit(api, db)

		// Initialize MissionsOrders routes
		mission_orders.MissionOrdersRouterInit(api, db)

		// Initialize Questions routes
		questions.QuestionRouterInit(api, db)

		// Initialize Tests routes
		tests.TestRouterInit(api, db)

		//  Initialize condidats routes
		candidats.CandidatRouterInit(api, db)

		//  Initialize projects routes
		projects.ProjectRouterInit(api, db)
	}
}
