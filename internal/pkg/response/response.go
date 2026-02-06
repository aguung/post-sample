package response

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	reqID, _ := c.Get("RequestID")
	c.JSON(code, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: reqID.(string),
	})
}

func Error(c *gin.Context, code int, message string, err interface{}) {
	reqID, _ := c.Get("RequestID")

	var errorResponse interface{} = err

	if e, ok := err.(error); ok {
		if e == io.EOF {
			message = "Request body cannot be empty"
			errorResponse = nil
		} else if _, ok := e.(validator.ValidationErrors); ok {
			message = "Validation error"
			errorResponse = FormatValidationError(e)
		} else {
			errorResponse = e.Error()
		}
	}

	c.JSON(code, Response{
		Success:   false,
		Message:   message,
		Error:     errorResponse,
		RequestID: reqID.(string),
	})
}
