/*

	Package utils provides utility functions and structures, primarily for handling JWT tokens, authorization, and extraction of user session information in a Gin web framework context.

	Structures:
	- UserSessionJWT: Structure for storing user session information extracted from JWT tokens.

	Functions:
	- GenerateToken(id, companyID uuid.UUID, roles []domains.RolesSessionJWT) string: Generates a JWT token with the provided user and role information.
	- AuthorizeJWT() gin.HandlerFunc: Gin middleware to authorize JWT tokens.
	- validateToken(token string) (*jwt.Token, error): Validates a JWT token and returns the token if valid.
	- ExtractJWTValues(ctx *gin.Context) domains.UserSessionJWT: Extracts user session values from the JWT token present in the request header.
	- extractToken(ctx *gin.Context) string: Extracts the JWT token from the Authorization header in the request.

	Dependencies:
	- "fmt": Standard Go package for formatted I/O.
	- "labs/constants": Custom package for application-specific constants.
	- "labs/domains": Custom package for application-specific domain structures.
	- "net/http": Standard Go package for HTTP protocols.
	- "github.com/gin-gonic/gin": Web framework for building APIs in Go.
	- "github.com/golang-jwt/jwt/v4": JWT implementation for Go.
	- "github.com/google/uuid": Package for UUID generation.
	- "github.com/sirupsen/logrus": Structured logger for Go.

	Usage:
	- Import this package to utilize the provided utility functions for handling JWT tokens and user session information in a Gin web framework context.

	Note:
	- The UserSessionJWT structure is used for storing user session information extracted from JWT tokens.
	- GenerateToken creates a JWT token with specified user and role information.
	- AuthorizeJWT is a Gin middleware for JWT token authorization.
	- validateToken is a utility function for validating JWT tokens.
	- ExtractJWTValues extracts user session information from the JWT token present in the request header.
	- extractToken is a utility function to extract the JWT token from the Authorization header in the request.

	Last update :
	01/02/2024 10:22

*/

package utils

import (
	"fmt"
	"labs/constants"
	"labs/domains"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GenerateToken generates a JWT token with the provided user and role information.
func GenerateToken(id, companyID uuid.UUID, roles []domains.RolesSessionJWT) string {

	// Retrieve JWT duration from environment variable
	duration, err := GetIntEnv("JWT_DURATION")
	if err != nil {
		logrus.Fatal("An error occurred when reading value from env. Error", err.Error())
	}

	// Retrieve JWT secret from environment variable
	secret, err := GetStringEnv("JWT_SECRET")
	if err != nil {
		logrus.Fatal("An error occurred when reading value from env. Error", err.Error())
	}

	// Set JWT claims including expiration time, issued at time, user ID, company ID, and roles
	claims := jwt.MapClaims{
		"exp":        time.Now().Add(time.Hour * time.Duration(duration)).Unix(),
		"iat":        time.Now().Unix(),
		"user_id":    id,
		"company_id": companyID,
		"roles":      roles,
	}

	// Create a new JWT token with specified claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the specified secret and handle potential errors
	auth, err := token.SignedString([]byte(secret))
	if err != nil {
		logrus.Fatal("An error occurred when reading value from env. Error", err.Error())
	}

	return auth
}

// AuthorizeJWT is a Gin middleware to authorize JWT tokens.
func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BearerSchema string = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			BuildErrorResponse(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED, Null())
			return
		}
		tokenString := authHeader[len(BearerSchema):]
		if token, err := validateToken(tokenString); err != nil {
			logrus.WithFields(logrus.Fields{
				"token": tokenString,
				"error": err.Error(),
			}).Fatal("Token validation error")
			BuildErrorResponse(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED, Null())
			return
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); !ok {
				BuildErrorResponse(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED, Null())
				return
			} else {
				if token.Valid {
					ctx.Set("user_id", claims["user_id"])
					ctx.Set("company_id", claims["company_id"])
					ctx.Set("roles", claims["roles"])
				} else {
					BuildErrorResponse(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED, Null())
					return
				}
			}
		}
	}
}

// validateToken validates a JWT token and returns the token if valid.
func validateToken(token string) (*jwt.Token, error) {
	// Use 'Parse' function to decode the token and validate its signature
	// The 2nd argument is a function that returns the secret key after checking if the signing method is HMAC
	// The returned key is used by 'Parse' to decode the token
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// If not HMAC, return a nil secret key and log an error
			logrus.WithFields(logrus.Fields{
				"unexpected_signing_method": token.Header["alg"],
			}).Fatal("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve JWT secret from environment variable
		secret, err := GetStringEnv("JWT_SECRET")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("An error occurred when reading value from env.")
		}

		// Return the secret key for validation
		return []byte(secret), nil
	})
}

// ExtractJWTValues extracts user session values from the JWT token present in the request header.
func ExtractJWTValues(ctx *gin.Context) domains.UserSessionJWT {

	// Initialize an empty UserSessionJWT struct to store extracted values
	session := domains.UserSessionJWT{}

	// Extract the JWT token string from the request
	tokenString := extractToken(ctx)

	// Parse the JWT token and validate its signature
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Log an error if the signing method is unexpected
			logrus.WithFields(logrus.Fields{
				"unexpected_signing_method": token.Header["alg"],
			}).Error("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve JWT secret from environment variable
		secret, err := GetStringEnv("JWT_SECRET")
		if err != nil {
			// Log an error if there's an issue retrieving the secret from the environment variable
			logrus.Fatal("An error occurred when reading value from env. Error ", err)
		}

		// Return the secret key for validation
		return []byte(secret), nil
	})

	// Check if the token is valid and contains the expected claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Parse and assign user ID
		userID, err := ParseUUID(claims["user_id"])
		if err != nil {
			logrus.Fatal("Error parsing user ID from JWT claims. Error ", err)
			return domains.UserSessionJWT{}
		}

		// Parse and assign company ID
		companyID, err := ParseUUID(claims["company_id"])
		if err != nil {
			logrus.WithError(err).Error("Error parsing company ID from JWT claims")
			return domains.UserSessionJWT{}
		}

		// Assign parsed values to the UserSessionJWT struct
		session.UserID = userID
		session.CompanyID = companyID

		// Extract roles as []interface{} and convert each to RolesSessionJWT
		rolesClaim, ok := claims["roles"].([]interface{})
		if !ok {
			logrus.Error("Error extracting roles from JWT claims")
			return domains.UserSessionJWT{}
		}

		// Loop through role claims and convert to RolesSessionJWT
		for _, roleClaim := range rolesClaim {
			roleData, ok := roleClaim.(map[string]interface{})
			if !ok {
				logrus.Error("Error extracting role data from JWT claims")
				return domains.UserSessionJWT{}
			}

			roleID, err := ParseUUID(roleData["id"])
			if err != nil {
				logrus.WithError(err).Error("Error parsing role ID from JWT claims")
				return domains.UserSessionJWT{}
			}

			roleName, ok := roleData["name"].(string)
			if !ok {
				logrus.Error("Error extracting role name from JWT claims")
				return domains.UserSessionJWT{}
			}

			companyID, err := ParseUUID(roleData["company_id"])
			if err != nil {
				logrus.WithError(err).Error("Error parsing role company ID from JWT claims")
				return domains.UserSessionJWT{}
			}

			// Append the RolesSessionJWT to the session.Roles slice
			session.Roles = append(session.Roles, domains.RolesSessionJWT{
				ID:        roleID,
				Name:      roleName,
				CompanyID: companyID,
			})
		}

		// Return the filled UserSessionJWT struct
		return session
	}

	// Return an empty UserSessionJWT if token validation fails
	return domains.UserSessionJWT{}
}

// extractToken extracts the JWT token from the Authorization header in the request.
func extractToken(ctx *gin.Context) string {
	// Retrieve the Authorization header from the request
	authorizationHeader := ctx.Request.Header["Authorization"]

	// Check if the Authorization header is present and has the expected format
	if len(authorizationHeader) > 0 {
		bearerToken := strings.Fields(authorizationHeader[0])[1]

		// Return the extracted token
		return bearerToken
	}

	// Return an empty string if the Authorization header is not present or doesn't have the expected format
	return ""
}
