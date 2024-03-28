package controller

import (
	"github.com/core-api/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func (ctl *Controller) CreateWordPress(c *gin.Context) {
	userResponse, code, err := ctl.svc.CreateWordPress()
	// userResponse := "controller for wordpress"
	if err != nil {
		response.ERROR(c, err, code)
	}
	response.JSON(c, userResponse, "Success", 0, 0)
}
