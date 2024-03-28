package controller

import (
	"github.com/core-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Router *gin.Engine
	svc    service.ServiceInterface
}

func NewController(svc service.ServiceInterface) *Controller {
	ctl := &Controller{}
	ctl.Router = gin.Default()
	ctl.svc = svc
	ctl.Routes()
	return ctl
}
