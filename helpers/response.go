package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONResponse represents the structure of a standardized API response
type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondSuccess sends a success response with data
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// RespondError sends an error response with a custom status code
func RespondError(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, JSONResponse{
		Status:  "error",
		Message: message,
		Error:   err.Error(),
	})
}
