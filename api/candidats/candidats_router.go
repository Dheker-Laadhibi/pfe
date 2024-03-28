package candidats

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CondidatRouterInit initializes the routes related to condidats.
func CandidatRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewCondidatRepository(db)

	// Private
	candidats := router.Group("/candidats")
	{

		// POST endpoint to create a new condidat
		candidats.POST("/:companyID", baseInstance.Createcandidate)

		// GET endpoint to retrieve all condidats for a specific company
		candidats.GET("/:companyID", baseInstance.ReadCandidats)

		// GET endpoint to retrieve a list of condidats for a specific company
		candidats.GET("/:companyID/list", baseInstance.ReadCondidatsList)

		// GET endpoint to retrieve the count of condidats for a specific company
		candidats.GET("/:companyID/count", baseInstance.ReadCandidatsCount)

		// GET endpoint to retrieve details of a specific condidat
		candidats.GET("/:companyID/:ID", baseInstance.Readcandidat)

		// PUT endpoint to update a specific condidat
		candidats.PUT("/:companyID/:ID", baseInstance.Updatecandidat)

		// DELETE endpoint to delete a specific condidat
		candidats.DELETE("/:companyID/:ID", baseInstance.DeleteCondidat)

		// signin endpoint to  a specific condidat
		candidats.POST("/signin", baseInstance.SigninCandidat)

	}
}
