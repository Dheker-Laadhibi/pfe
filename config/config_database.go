/*

	Package config provides functionality for configuring and connecting to the database.

	The ConnectToDB function initializes and returns a connection to the PostgreSQL database.

	Functions:
	- ConnectToDB(): Initializes and returns a connection to the PostgreSQL database.
	- Flags(database *gorm.DB): Parses command line flags, allowing the creation of a root user.

	Dependencies:
	- "labs/utils": Custom package for utility functions.
	- "log": Standard Go logging package.
	- "os": Standard Go package for interacting with the operating system.
	- "time": Standard Go package for handling time.
	- "github.com/sirupsen/logrus": Structured logger for Go (used for enhanced logging).
	- "gorm.io/driver/postgres": GORM PostgreSQL driver.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.
	- "gorm.io/gorm/logger": GORM logger package for customizing database logging.

	Environment Variables:
	- DB_DSN: Database connection string.

	Usage:
	- Call ConnectToDB() to establish a connection to the PostgreSQL database.
	- Call Flags(database *gorm.DB) to parse command line flags for creating a root user.

	Note:
	- The ConnectToDB function sets the client encoding to 'UTF8' for the database.

	Last update :
	14/02/2024 10:22

*/

package config

import (
	"flag"
	"log"
	"os"
	"time"

	"labs/utils"

	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectToDB initializes and returns a connection to the PostgreSQL database.
func ConnectToDB() *gorm.DB {

	// Create a new logger for database queries
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second * 10,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)

	// Get the database connection string from the environment
	dns, err := utils.GetStringEnv("DB_DSN")
	if err != nil {
		logrus.Fatal("Failed to load DNS from env file: ", err)
	}

	// Open a connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{Logger: newLogger})
	if err != nil {
		logrus.Fatal("Failed to connect to the database: ", err)
	}

	// Set the client encoding for the database to UTF8
	if err := db.Exec("SET client_encoding = 'UTF8'").Error; err != nil {
		logrus.Fatal("Failed to set database client encoding: ", err)
	}

	// Return the initialized database connection
	return db
}

// Flags parses command line flags, allowing the creation of a root user.
func Flags(database *gorm.DB) {

	// Define a flag to create a root user
	root := flag.Bool("root", false, "Create Root")
	flag.Parse()

	// Check if the root flag is set
	if *root {

		// If root flag is set, create root role, user, and company using the Seeder function
		if err := Seeder(database); err != nil {
			logrus.Fatal("Failed to create root user: ", err.Error())
		}
	}
}
