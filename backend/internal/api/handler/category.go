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

type CategoryHandler struct {
	categoryService service.ICategoryService
}

func NewCategoryHandler(categoryService service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategories 获取所有分类
func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(ctx.Request.Context())
	if err != nil {
		logger.Log.Warn("获取分类列表失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取分类列表失败")
		return
	}
	result.Success(ctx, "获取分类列表成功", categories)
}

// GetCategoryByID 根据 ID 获取分类详情
func (h *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Warn("获取分类详情失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的分类 ID")
		return
	}

	category, err := h.categoryService.GetCategoryByID(ctx.Request.Context(), uint64(id))
	if err != nil {
		logger.Log.Warn("获取分类详情失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取分类详情失败")
		return
	}
	result.Success(ctx, "获取分类详情成功", category)
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		logger.Log.Warn("创建分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.categoryService.CreateCategory(ctx.Request.Context(), &category); err != nil {
		logger.Log.Warn("创建分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "创建分类失败")
		return
	}
	result.Success(ctx, "创建分类成功", category)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Warn("更新分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的分类 ID")
		return
	}
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		logger.Log.Warn("更新分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}
	category.ID = uint64(id)
	err = h.categoryService.UpdateCategory(ctx.Request.Context(), &category)
	if err != nil {
		logger.Log.Warn("更新分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "更新分类失败")
		return
	}
	result.Success(ctx, "更新分类成功", category)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Warn("删除分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的分类 ID")
		return
	}

	if err := h.categoryService.DeleteCategory(ctx.Request.Context(), uint64(id)); err != nil {
		logger.Log.Warn("删除分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "删除分类失败")
		return
	}
	result.Success(ctx, "删除分类成功", nil)
}
