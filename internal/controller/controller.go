package controller

import "github.com/gin-gonic/gin"

type Controller struct {
	Router *gin.Engine
}

func NewController() *Controller {
	ctl := &Controller{}
	ctl.Router = gin.Default()
	ctl.Routes()
	return ctl
}
