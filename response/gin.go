//go:build gin
// +build gin

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinJSON sends a JSON response with the given status code and data
func GinJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Code: statusCode,
		Data: data,
	})
}

// GinSuccess sends a successful JSON response (200 OK)
func GinSuccess(c *gin.Context, data interface{}) {
	GinJSON(c, http.StatusOK, data)
}

// GinCreated sends a created JSON response (201 Created)
func GinCreated(c *gin.Context, data interface{}) {
	GinJSON(c, http.StatusCreated, data)
}

// GinNoContent sends a no content response (204 No Content)
func GinNoContent(c *gin.Context) {
	c.AbortWithStatus(http.StatusNoContent)
}

// GinError sends an error JSON response
func GinError(c *gin.Context, statusCode int, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(statusCode, Response{
		Code:  statusCode,
		Error: errMsg,
	})
}

// GinBadRequest sends a bad request error (400 Bad Request)
func GinBadRequest(c *gin.Context, err error) {
	GinError(c, http.StatusBadRequest, err)
}

// GinUnauthorized sends an unauthorized error (401 Unauthorized)
func GinUnauthorized(c *gin.Context, err error) {
	GinError(c, http.StatusUnauthorized, err)
}

// GinForbidden sends a forbidden error (403 Forbidden)
func GinForbidden(c *gin.Context, err error) {
	GinError(c, http.StatusForbidden, err)
}

// GinNotFound sends a not found error (404 Not Found)
func GinNotFound(c *gin.Context, err error) {
	GinError(c, http.StatusNotFound, err)
}

// GinInternalServerError sends an internal server error (500 Internal Server Error)
func GinInternalServerError(c *gin.Context, err error) {
	GinError(c, http.StatusInternalServerError, err)
}

// GinStatusCode sends a response with a custom status code and message
func GinStatusCode(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Code:  statusCode,
		Error: message,
	})
}

// GinPaginated sends a paginated JSON response
func GinPaginated(c *gin.Context, data interface{}, page, pageSize, total int) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Code:     http.StatusOK,
		Data:     data,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Code     int         `json:"code"`
	Data     interface{} `json:"data,omitempty"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
}
