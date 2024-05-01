package companies

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CompanyRouterInit initializes the routes related to companies.
func CompanyRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewCompanyRepository(db)

	// Private
	companies := router.Group("/companies")
	{ 
		
		// create a new company 
		companies.POST("", baseInstance.CreateCompany)
		

		// GET endpoint to retrieve all companies
		companies.GET("", baseInstance.ReadCompanies)

		// GET endpoint to retrieve details of a specific company
		companies.GET("/:ID", baseInstance.ReadCompany)

		// PUT endpoint to update details of a specific company
		companies.PUT("/:ID", baseInstance.UpdateCompany)

		// DELETE endpoint to delete a specific company
		companies.DELETE("/:ID", baseInstance.DeleteCompany)
	}
}
