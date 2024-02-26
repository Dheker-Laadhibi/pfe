package constants

// Constant API Responses
const (
	SUCCESS         = "Success"
	CREATED         = "Created"
	DATA_NOT_FOUND  = "Data Not Found"
	UNKNOWN_ERROR   = "Unknown Error"
	INVALID_REQUEST = "Invalid Request"
	UNAUTHORIZED    = "Unauthorized"
	SERVER_ERROR    = "Server Error"
)

// Constant Regex
const (
	EMPTY_REGEX = `^( | )+$|^$`
	EMAIL_REGEX = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

// Constant Values
const (
	DEFAULT_ROLE             = "Manager"
	DEFAULT_PAGE_PAGINATION  = 1
	DEFAULT_LIMIT_PAGINATION = 10
)
