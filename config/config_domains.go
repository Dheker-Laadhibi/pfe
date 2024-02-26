/*
	Package config provides functionality related to application configuration, including initialization and seeding of data.

	The functions in this package handle the initialization of essential entities such as companies, users, roles, and their relationships.

	Functions:
	- parseValues(row string) []string: Divides a CSV string of key:value pairs into slices.
	- initCompany() (*domains.Companies, error): Initializes a new company based on environment variables.
	- initUser() (*domains.Users, error): Initializes a new user based on environment variables.
	- initRole() (*domains.Roles, error): Initializes a new role based on environment variables.
	- initUsersRoles(userID, roleID, companyID uuid.UUID) (*domains.UsersRoles, error): Initializes the relationship between a user and a role.
	- Seeder(db *gorm.DB) error: Seeds essential data into the database during application initialization.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "github.com/sirupsen/logrus": Structured logger for Go.
	- "golang.org/x/crypto/bcrypt": Package for hashing and verifying passwords.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.
	- "labs/domains": Package containing domain models and generic CRUD operations.

	Usage:
	- Import this package to utilize the configuration and seeding functions.

	Note:
	- The functions make use of environment variables for configuration.
	- The Seeder function initializes and seeds essential data into the database.

	Last update:
	14/02/2024 23:19

*/

package config

import (
	"errors"
	"labs/domains"
	"labs/utils"
	"log"
	"strconv"
	"strings"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// parseValues divides a CSV string of key:value pairs into slices
func parseValues(row string) []string {
	return strings.Split(row, ",")
}

// Function to initialize a new company based on environment variables
func initCompany() (*domains.Companies, error) {

	env_company, err := utils.GetStringEnv("ROOT_COMPANY")
	if err != nil || env_company == "" {
		logrus.Error("Happened error when init new company. Error", err.Error())
		return &domains.Companies{}, errors.New("Happened error when init new company " + err.Error())
	}

	company := &domains.Companies{}
	root_values := parseValues(env_company)

	company.ID = uuid.New()
	company.Name = root_values[0]

	return company, nil
}

// Function to initialize a new user based on environment variables
func initUser() (*domains.Users, error) {

	env_user, err := utils.GetStringEnv("ROOT_USER")
	if err != nil || env_user == "" {
		logrus.Error("Happened error when init new user. Error: ", err.Error())
		return &domains.Users{}, errors.New("Happened error when init new user " + err.Error())
	}

	user := &domains.Users{}
	root_values := parseValues(env_user)
	log.Println(root_values[3])
	hash, _ := bcrypt.GenerateFromPassword([]byte(root_values[3]), bcrypt.DefaultCost)

	user.ID = uuid.New()
	user.Firstname = root_values[0]
	user.Lastname = root_values[1]
	user.Email = root_values[2]
	user.Password = string(hash)
	user.Country = root_values[4]
	user.Status, _ = strconv.ParseBool(root_values[5])

	return user, nil
}

// Function to initialize a new role based on environment variables
func initRole() (*domains.Roles, error) {

	env_role, err := utils.GetStringEnv("ROOT_ROLE")
	if err != nil || env_role == "" {
		logrus.Error("Happened error when init new role. Error: ", err.Error())
		return &domains.Roles{}, errors.New("Happened error when init new role " + err.Error())
	}

	role := &domains.Roles{}
	root_values := parseValues(env_role)

	role.ID = uuid.New()
	role.Name = root_values[0]

	return role, nil
}

// Function to initialize the relationship between a user and a role
func initUsersRoles(user_id, role_id, company_id uuid.UUID) (*domains.UsersRoles, error) {

	user_role := &domains.UsersRoles{}
	user_role.UserID = user_id
	user_role.RoleID = role_id
	user_role.CompanyID = company_id

	return user_role, nil
}

// Function to seed essential data into the database during application initialization
func Seeder(db *gorm.DB) error {

	company, err := initCompany()
	if err != nil {
		return err
	}
	err = domains.Create(db, &company)
	if err != nil {
		return err
	}

	user, err := initUser()
	if err != nil {
		return err
	}
	user.CompanyID = company.ID
	user.CreatedByUserID = user.ID
	err = domains.Create(db, &user)
	if err != nil {
		return err
	}

	role, err := initRole()
	if err != nil {
		return err
	}
	role.OwningCompanyID = company.ID
	role.CreatedByUserID = user.ID
	err = domains.Create(db, &role)
	if err != nil {
		return err
	}

	userRole, err := initUsersRoles(user.ID, role.ID, company.ID)
	if err != nil {
		return err
	}
	err = domains.Create(db, &userRole)
	if err != nil {
		return err
	}

	return nil
}
