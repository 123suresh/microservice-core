package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success int         `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count,omitempty"`
	Limit   int         `json:"limit,omitempty"`
}

func ERROR(c *gin.Context, err error, code int) {
	c.JSON(code, Response{
		Message: err.Error(),
		Success: 0,
	})
}

func JSON(c *gin.Context, data interface{}, message string, limit int, count int) {
	c.JSON(http.StatusOK, Response{
		Success: 1,
		Message: message,
		Data:    data,
		Limit:   limit,
		Count:   count,
	})
}
