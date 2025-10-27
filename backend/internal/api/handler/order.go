package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/api/request"
	"github.com/wangn-tech/bookstore-go/internal/api/response"
	"github.com/wangn-tech/bookstore-go/internal/service"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// GetOrderStatistics 获取订单统计数据
func (o *OrderHandler) GetOrderStatistics(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("GetOrderStatistics: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	// stats 在 repository 层封装为 map[string]any 返回, 正常应该定义 VO 结构体
	stats, err := o.orderService.GetOrderStatistics(ctx, userID.(uint64))
	if err != nil {
		logger.Log.Error("GetOrderStatistics: 获取订单统计数据失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取订单统计数据失败")
		return
	}

	result.Success(ctx, "获取订单统计数据成功", stats)
}

// PayOrder 支付订单
func (o *OrderHandler) PayOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Warn("PayOrder: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的订单 ID")
		return
	}
	err = o.orderService.PayOrder(ctx, uint64(id))
	if err != nil {
		logger.Log.Error("PayOrder: 支付订单失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "支付订单失败")
		return
	}
	result.Success(ctx, "支付订单成功", nil)
}

// GetUserOrders 获取用户的订单列表
func (o *OrderHandler) GetUserOrders(ctx *gin.Context) {
	var req request.OrdersPageDTO
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Log.Warn("GetUserOrders: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}
	// 从 context 中获取 userID
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("GetUserOrders: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	result.PageVerify(&req.Page, &req.PageSize)
	// 调用 service 层获取订单列表 (分页查询)
	pageResult, err := o.orderService.GetUserOrders(ctx, userID.(uint64), &req)
	if err != nil {
		logger.Log.Error("GetUserOrders: 获取订单列表失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取订单列表失败")
		return
	}
	// 封装 VO 对象并返回
	result.Success(ctx, "获取订单列表成功", &response.OrdersPageVO{
		Orders:    pageResult.Records,
		Total:     pageResult.Total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: (pageResult.Total + int64(req.PageSize) - 1) / int64(req.PageSize),
	})
}

// GetOrderByID 根据 ID 获取订单详情
func (o *OrderHandler) GetOrderByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Warn("GetOrderByID: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的订单 ID")
		return
	}
	order, err := o.orderService.GetOrderByID(ctx, uint64(id))
	if err != nil {
		logger.Log.Error("GetOrderByID: 获取订单失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取订单失败")
		return
	}
	result.Success(ctx, "获取订单成功", order)
}

// CreateOrder 创建订单
func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req request.CreateOrderDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("CreateOrder: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}
	// 从 context 中获取 userID
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("CreateOrder: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	req.UserID = userID.(uint64)
	order, err := o.orderService.CreateOrder(ctx, &req)
	if err != nil {
		logger.Log.Error("CreateOrder: 创建订单失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "创建订单失败")
		return
	}
	result.Success(ctx, "创建订单成功", order)
}
