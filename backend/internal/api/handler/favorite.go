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

type FavoriteHandler struct {
	favoriteService service.IFavoriteService
}

func NewFavoriteHandler(favoriteService service.IFavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

func (h *FavoriteHandler) AddFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("AddFavorite: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		logger.Log.Warn("AddFavorite: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的书籍 ID")
		return
	}
	err = h.favoriteService.AddFavorite(ctx, userID.(uint64), bookID)
	if err != nil {
		logger.Log.Error("AddFavorite: 添加收藏失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "添加收藏失败")
		return
	}
	result.Success(ctx, "添加收藏成功", nil)
}

// RemoveFavorite 删除收藏
func (h *FavoriteHandler) RemoveFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("RemoveFavorite: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		logger.Log.Warn("RemoveFavorite: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的书籍 ID")
		return
	}
	err = h.favoriteService.RemoveFavorite(ctx, userID.(uint64), bookID)
	if err != nil {
		logger.Log.Error("RemoveFavorite: 删除收藏失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "删除收藏失败")
		return
	}
	result.Success(ctx, "取消收藏成功", nil)
}

// CheckFavorite 检查是否已收藏
func (h *FavoriteHandler) CheckFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("CheckFavorite: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	bookIDStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	// bookIDInt, err := strconv.Atoi(bookIDStr)
	// bookID := uint64(bookIDInt)
	if err != nil {
		logger.Log.Warn("CheckFavorite: 请求参数错误", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的书籍 ID")
		return
	}
	isFavorited, err := h.favoriteService.IsFavorited(ctx, userID.(uint64), bookID)
	if err != nil {
		logger.Log.Error("CheckFavorite: 检查收藏状态失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "检查收藏失败")
		return
	}
	result.Success(ctx, "检查收藏成功", isFavorited)
}

// GetUserFavorites 获取用户所有收藏的书籍
func (h *FavoriteHandler) GetUserFavorites(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("GetUserFavorites: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	var pageReq request.FavoritePageDTO
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		logger.Log.Warn("GetUserFavorites: 请求参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}
	result.PageVerify(&pageReq.Page, &pageReq.PageSize)
	timeFilter := ctx.DefaultQuery("time_filter", "all")
	pageResult, err := h.favoriteService.GetUserFavorites(ctx, userID.(uint64), &pageReq, timeFilter)
	if err != nil {
		logger.Log.Error("GetUserFavorites: 获取用户收藏列表失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取用户收藏列表失败")
		return
	}
	result.Success(ctx, "获取用户收藏列表成功", &response.FavoritesPageVO{
		Favorites:   pageResult.Records,
		Total:       pageResult.Total,
		CurrentPage: pageReq.Page,
		// PageSizes:   pageReq.PageSize,
		TotalPages: (pageResult.Total + int64(pageReq.PageSize) - 1) / int64(pageReq.PageSize),
	})
}

// GetUserFavoriteCount 获取用户收藏数量
func (h *FavoriteHandler) GetUserFavoriteCount(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.Log.Warn("GetUserFavoriteCount: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	count, err := h.favoriteService.GetUserFavoriteCount(ctx, userID.(uint64))
	if err != nil {
		logger.Log.Error("GetUserFavoriteCount: 获取用户收藏数量失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取用户收藏数量失败")
		return
	}

	result.Success(ctx, "获取用户收藏数量成功", count)
}
