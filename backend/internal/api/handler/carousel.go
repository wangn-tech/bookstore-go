package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/service"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
	"go.uber.org/zap"
)

type CarouselHandler struct {
	carouselService service.ICarouselService
}

func NewCarouselHandler(carouselService service.ICarouselService) *CarouselHandler {
	return &CarouselHandler{
		carouselService: carouselService,
	}
}

// GetActiveCarousels 获取所有激活的轮播图
func (h *CarouselHandler) GetActiveCarousels(ctx *gin.Context) {
	carousels, err := h.carouselService.GetActiveCarousels(ctx.Request.Context())
	if err != nil {
		logger.Log.Warn("GetActiveCarousels: Failed to get active carousels", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取轮播图失败")
		return
	}
	result.Success(ctx, "获取轮播图成功", carousels)
}

// GetCarouselByID 获取指定ID的轮播图
func (h *CarouselHandler) GetCarouselByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("GetCarouselByID: Invalid carousel ID", zap.String("id", idStr), zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的轮播图ID")
		return
	}
	carousel, err := h.carouselService.GetCarouselByID(ctx.Request.Context(), id)
	if err != nil {
		logger.Log.Warn("GetCarouselByID: Failed to get carousel", zap.Uint64("id", id), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取轮播图失败")
		return
	}
	result.Success(ctx, "获取轮播图成功", carousel)
}

// CreateCarousel 创建轮播图
func (h *CarouselHandler) CreateCarousel(ctx *gin.Context) {
	var carousel model.Carousel
	if err := ctx.ShouldBindJSON(&carousel); err != nil {
		logger.Log.Warn("CreateCarousel: Invalid request body", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "请求参数无效")
		return
	}
	if err := h.carouselService.CreateCarousel(ctx.Request.Context(), &carousel); err != nil {
		logger.Log.Warn("CreateCarousel: Failed to create carousel", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "创建轮播图失败")
		return
	}
	result.Success(ctx, "创建轮播图成功", carousel)
}

// UpdateCarousel 更新轮播图
func (h *CarouselHandler) UpdateCarousel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("UpdateCarousel: Invalid carousel ID", zap.String("id", idStr), zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的轮播图ID")
		return
	}
	var carousel model.Carousel
	if err := ctx.ShouldBindJSON(&carousel); err != nil {
		logger.Log.Warn("UpdateCarousel: Invalid request body", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "请求参数无效")
		return
	}
	carousel.ID = id
	err = h.carouselService.UpdateCarousel(ctx.Request.Context(), &carousel)
	if err != nil {
		logger.Log.Warn("UpdateCarousel: Failed to update carousel", zap.Uint64("id", id), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "更新轮播图失败")
		return
	}
	result.Success(ctx, "更新轮播图成功", carousel)
}

// DeleteCarousel 删除轮播图
func (h *CarouselHandler) DeleteCarousel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("DeleteCarousel: Invalid carousel ID", zap.String("id", idStr), zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的轮播图ID")
		return
	}
	carousel, err := h.carouselService.DeleteCarousel(ctx.Request.Context(), id)
	if err != nil {
		logger.Log.Warn("DeleteCarousel: Failed to get carousel", zap.Uint64("id", id), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "删除轮播图失败")
		return
	}
	result.Success(ctx, "删除轮播图成功", carousel)