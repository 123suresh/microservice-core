package controller

import (
	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Routes() {
	ctl.Router.Use(CORSMiddleware())
	ctl.Router.POST("/wordpress", ctl.CreateWordPress)
	ctl.Router.GET("/wordpress", ctl.GetWordPress)
	ctl.Router.DELETE("/wordpress", ctl.DeleteWordpress)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-agent-code")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
