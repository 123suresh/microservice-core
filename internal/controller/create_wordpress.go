package controller

import (
	"github.com/core-api/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func (ctl *Controller) CreateWordPress(c *gin.Context) {
	// userResponse, code, err := "controller for wordpress"
	// userResponse := "controller for wordpress"
	// if err != nil {
	// 	response.ERROR(c, err, code)
	// }
	response.JSON(c, "controller for wordpress", "Success", 0, 0)
}
