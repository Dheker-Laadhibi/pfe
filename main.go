package main

import (
	"labs/api"
	"labs/config"
	"labs/utils"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	godotenv.Load()
	config.InitLog()
}

// Main is the entry point of the application.
// @title           			Project Name
// @version         			0.1
// @description					A basic template for managing your project using Golang REST APIs.
// @contact.name    			Hedi Manai
// @contact.email				h.manai@tunisiancloud.com
// @host     					localhost:8080
// @BasePath 					/api
// @schemes						http
// @securityDefinitions.apikey	ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// Retrieve the port number from the environment variables
	port, err := utils.GetStringEnv("PORT")
	if err != nil {
		logrus.Fatal("Failed to load PORT from env file: ", err)
	}

	// Initialize a router using the utility function InitRouter
	router := utils.InitRouter()

	// Connect to the PostgreSQL database and obtain a database connection
	database := config.ConnectToDB()

	// Initialize API routes using the router and the database connection
	api.RoutesApiInit(router, database)

	// Parse command line flags, allowing the creation of a root user
	config.Flags(database)

	// Run the web server on the specified port
	router.Run(":" + port)
}
