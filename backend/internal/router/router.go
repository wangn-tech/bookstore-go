package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/middlerware"
	"github.com/wangn-tech/bookstore-go/internal/router/bookstore"
)

type RouteGroup struct {
	bookstore.UserRouter
	bookstore.BookRouter
	bookstore.CategoryRouter
	bookstore.OrderRouter
	bookstore.FavoriteRouter
	bookstore.CarouselRouter
}

var AllRouter = new(RouteGroup)

func InitRouter(r *gin.Engine) {
	// r.Use(func(c *gin.Context) {
	// 	c.Header("Content-Type", "application/json; charset=utf-8")
	// 	c.Next()
	// })
	// 全局中间件 CORS
	r.Use(middlerware.CORS())

	// 健康检查 (ping, pong)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 业务路由
	v1 := r.Group("/api/v1")
	{
		AllRouter.InitUserRouter(v1)
		AllRouter.InitBookRouter(v1)
		AllRouter.InitCategoryRouter(v1)
		AllRouter.InitOrderRouter(v1)
		AllRouter.InitFavoriteRouter(v1)
		AllRouter.InitCarouselRouter(v1)

	}
}
