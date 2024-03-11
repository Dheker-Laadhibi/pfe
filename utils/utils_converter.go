/*

	package utils provides utility functions for handling common tasks such as UUID parsing,
	environment variable retrieval, password comparison, and more.

	UUID Parsing:
	The ParseUUID function parses an input ID into a UUID. It supports both UUID and string
	inputs. If the input is already a UUID, it is returned as is. If it is a string, the
	function attempts to parse it into a UUID using the github.com/google/uuid package.

	Environment Variable Retrieval:
	Several functions are provided for retrieving environment variables of different types.

	- GetStringEnv: Retrieves a string environment variable and trims any leading or trailing whitespaces.

	- GetBoolEnv: Retrieves a boolean environment variable and parses it using strconv.ParseBool.

	- GetIntEnv: Retrieves an integer environment variable and converts it using strconv.Atoi.

	- GetArrayEnv: Retrieves an environment variable as a string and splits it into an array of strings using a specified separator.

	Password Comparison:
	The ComparePassword function uses bcrypt.CompareHashAndPassword to compare a hashed password
	stored in the database with a plaintext password. It returns true if the passwords match.

	Note:
	It is crucial to handle errors returned by these functions appropriately in your application
	to ensure robust error handling and improve overall reliability.

	Repeat similar usage for other functions based on your requirements.

	Last update :
	01/02/2024 10:22

*/

package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ParseUUID(id interface{}) (uuid.UUID, error) {

	// Switch on the type of the 'id' variable
	switch id := id.(type) {

	// If 'id' is already a UUID, return it as is
	case uuid.UUID:
		return id, nil

	// If 'id' is a string, attempt to parse it into a UUID
	case string:

		// Use the uuid.Parse function to parse the string into a UUID
		parsedUUID, err := uuid.Parse(id)

		// If there's an error during parsing, return the nil UUID and the error
		if err != nil {
			return uuid.Nil, err
		}

		// If parsing is successful, return the parsed UUID
		return parsedUUID, nil

	// If 'id' is of any other type, return the nil UUID and an error indicating it's not a valid UUID
	default:
		return uuid.Nil, errors.New("is not a valid UUID")
	}
}
// variable définie dans le projet  qui stocke des informations statiques comme port et root company  
func GetStringEnv(name string) (string, error) {

	// Retrieve the value of the environment variable with the specified name
	value := os.Getenv(name)

	// If the value is an empty string, return an error indicating the environment variable is not set
	if value == "" {
		return "", errors.New("environment variable is not set")
	}
//Suppression des espaces vides 
	// Trim any leading or trailing whitespaces from the environment variable value
	return strings.TrimSpace(value), nil
}

func GetBoolEnv(name string) (bool, error) {

	// Retrieve the value of the environment variable as a string using GetStringEnv
	value, err := GetStringEnv(name)

	// If there's an error during the retrieval, return false and the error
	if err != nil {
		return false, err
	}

	// Attempt to parse the string value as a boolean
	parsedBool, err := strconv.ParseBool(value)

	// If there's an error during parsing, return false and the parsing error
	if err != nil {
		return false, err
	}

	// Return the parsed boolean value and a nil error if parsing is successful
	return parsedBool, nil
}

func GetIntEnv(name string) (int, error) {

	// Retrieve the value of the environment variable as a string using GetStringEnv
	value, err := GetStringEnv(name)

	// If there's an error during the retrieval, return 0 and the error
	if err != nil {
		return 0, err
	}

	// Attempt to convert the string value to an integer
	intValue, err := strconv.Atoi(value)

	// If there's an error during conversion, return 0 and an error indicating the type mismatch
	if err != nil {
		return 0, errors.New("the environment variable is not of type int")
	}

	// Return the parsed integer value and a nil error if conversion is successful
	return intValue, nil
}

func GetArrayEnv(name string, separator string) ([]string, error) {

	// Retrieve the value of the environment variable as a string using GetStringEnv
	value, err := GetStringEnv(name)

	// If there's an error during the retrieval, return nil and the error
	if err != nil {
		return nil, err
	}

	// Split the string value into an array of strings using the specified separator
	arrayValue := strings.Split(value, separator)

	// Return the resulting array and a nil error
	return arrayValue, nil
}

func ComparePassword(dbPass, pass string) bool {

	// Use bcrypt.CompareHashAndPassword to compare the hashed password from the database with the plaintext password
	// The function returns nil if the passwords match, indicating success
	return bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass)) == nil
}
