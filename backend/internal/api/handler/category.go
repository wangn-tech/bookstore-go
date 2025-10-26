package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
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
		logger.Log.Error("获取分类列表失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取分类列表失败")
		return
	}
	result.Success(ctx, "获取分类列表成功", categories)
}

// GetCategoryByID 根据 ID 获取分类详情
func (h *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Error("获取分类详情失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "无效的分类 ID")
		return
	}

	category, err := h.categoryService.GetCategoryByID(ctx.Request.Context(), uint64(id))
	if err != nil {
		logger.Log.Error("获取分类详情失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取分类详情失败")
		return
	}
	result.Success(ctx, "获取分类详情成功", category)
}

func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	// Implementation for creating a category
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	// Implementation for updating a category
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	// Implementation for deleting a category
}
