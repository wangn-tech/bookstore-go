package bookstore

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/api/handler"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/middlerware"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/service"
)

type OrderRouter struct {
	orderService service.IOrderService
}

func (o *OrderRouter) InitOrderRouter(router *gin.RouterGroup) {
	o.orderService = service.NewOrderService(
		repository.NewOrderDao(database.DB),
		repository.NewBookDao(database.DB),
	)

	orderHandler := handler.NewOrderHandler(o.orderService)

	orderGroup := router.Group("/order")
	orderGroup.Use(middlerware.JWTAuth())
	{
		orderGroup.POST("/create", orderHandler.CreateOrder)           // 创建订单
		orderGroup.GET("/:id", orderHandler.GetOrderByID)              // 根据 ID 获取订单详情
		orderGroup.GET("/list", orderHandler.GetUserOrders)            // 获取用户订单列表
		orderGroup.POST("/:id/pay", orderHandler.PayOrder)             // 支付订单
		orderGroup.GET("/statistics", orderHandler.GetOrderStatistics) // 获取订单统计数据
	}
}
