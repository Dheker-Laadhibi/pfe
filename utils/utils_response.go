/*

	Package utils provides utility functions and structures, particularly for handling API responses in a Gin web framework context.

	Structures:
	- ApiResponses: Structure for standard API responses.
		- ResponseKey (string): Key indicating the response status or type.
		- Data (interface{}): Generic data payload associated with the response.

	Functions:
	- APIResponse(code string, data interface{}) ApiResponses: Creates and returns an ApiResponses instance with the specified response key and data.
	- BuildErrorResponse(ctx *gin.Context, status int, code string, data interface{}): Builds and sends a standardized error response with the specified code and data in the context of a Gin web framework.
	- BuildResponse(ctx *gin.Context, status int, code string, data interface{}): Builds and sends a standardized API response with the specified code and data in the context of a Gin web framework.
	- Null() interface{}: Returns a null interface.
	- ResponseLimitPagination() []uint: Returns pagination choices for response limits.

	Dependencies:
	- "net/http": Standard Go package for HTTP protocols.
	- "github.com/gin-gonic/gin": Web framework for building APIs in Go.

	Usage:
	- Import this package to utilize the provided utility functions for building API responses and handling errors in a Gin web framework context.

	Note:
	- The ApiResponses structure represents a standardized format for API responses.
	- APIResponse is used to create a response object with a specified code and data.
	- BuildErrorResponse is a utility for sending standardized error responses in the context of a Gin web framework.
	- BuildResponse is a utility for sending standardized API responses in the context of a Gin web framework.
	- Null is a utility function that returns a null interface.
	- ResponseLimitPagination provides choices for response limits in pagination.

	Last update :
	01/02/2024 10:22

*/

package utils

import (
	"github.com/gin-gonic/gin"
)

// @Description Generic API response
type ApiResponses struct {
	ResponseKey string      `json:"responseKey"`
	Data        interface{} `json:"data"`
} //@name ApiResponse

// APIResponse creates and returns an ApiResponses instance with the specified response key and data.
func APIResponse(code string, data interface{}) ApiResponses {
	return ApiResponses{
		ResponseKey: code,
		Data:        data,
	}
}

// BuildErrorResponse builds and sends a standardized error response with the specified code and data in the context of a Gin web framework.
func BuildErrorResponse(ctx *gin.Context, status int, code string, data interface{}) {
	ctx.JSON(status, APIResponse(code, data))
	ctx.Abort()
}

// BuildResponse builds and sends a standardized API response with the specified code and data in the context of a Gin web framework.
func BuildResponse(ctx *gin.Context, status int, code string, data interface{}) {
	ctx.JSON(status, APIResponse(code, data))
}

// Null returns a null interface.
func Null() interface{} {
	return nil
}

// ResponseLimitPagination returns pagination choices for response limits.
func ResponseLimitPagination() []uint {
	return []uint{5, 10, 20, 50}
}
