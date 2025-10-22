package router

import "github.com/gin-gonic/gin"

type RouteGroup struct {
	//bookstore.
}

var AllRouter = new(RouteGroup)

func InitRouter(r *gin.Engine) {

	// 测试路由ping, pong
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
