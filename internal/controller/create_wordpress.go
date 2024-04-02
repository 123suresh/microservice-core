package controller

import (
	"net/http"

	"github.com/core-api/internal/model"
	"github.com/core-api/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (ctl *Controller) CreateWordPress(c *gin.Context) {
	wordPressReq := &model.WordPressRequest{}
	err := c.ShouldBindJSON(&wordPressReq)
	if err != nil {
		logrus.Error("json bind error :: ", err)
		response.ERROR(c, err, http.StatusBadRequest)
		return
	}
	userResponse, code, err := ctl.svc.CreateWordPress(wordPressReq)
	// userResponse := "controller for wordpress"
	if err != nil {
		response.ERROR(c, err, code)
	}
	response.JSON(c, userResponse, "Success", 0, 0)
}
