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

type BookHandler struct {
	bookService service.IBookService
}

func NewBookHandler(bookService service.IBookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// GetBookList 获取书籍列表，支持分页
func (b *BookHandler) GetBookList(ctx *gin.Context) {
	// page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	// pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "12"))
	var pageReq request.BooksPageDTO
	// 绑定查询参数 (form)，如果绑定失败则使用默认值
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		// 对于分页查询，即使参数有误，也应提供默认查询结果，而不是直接报错
		logger.Log.Warn("PageQuery: 查询参数绑定失败，使用默认值", zap.Error(err))
		pageReq.Page = 1
		pageReq.PageSize = 10
	}

	pageResult, err := b.bookService.GetBooksByPage(ctx.Request.Context(), &pageReq)
	if err != nil {
		// Handle error (omitted for brevity)
		logger.Log.Error("GetBookList: 获取书籍列表失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取书籍列表失败")
		return
	}

	result.Success(ctx, "获取书籍列表成功", &response.BooksPageVO{
		Books:     pageResult.Records,
		Total:     pageResult.Total,
		Page:      pageReq.Page,
		PageSize:  pageReq.PageSize,
		TotalPage: (pageResult.Total + int64(pageReq.PageSize) - 1) / int64(pageReq.PageSize),
	})
}

// GetHotBooks 获取热销图书
func (b *BookHandler) GetHotBooks(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

	books, err := b.bookService.GetHotBooks(ctx.Request.Context(), limit)
	if err != nil {
		logger.Log.Error("GetHotBooks: 获取热销图书失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取热销图书失败")
		return
	}

	result.Success(ctx, "获取热销图书成功", books)
}

// GetNewBooks 获取新书
func (b *BookHandler) GetNewBooks(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	books, err := b.bookService.GetNewBooks(ctx.Request.Context(), limit)
	if err != nil {
		logger.Log.Error("GetNewBooks: 获取新书失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取新书失败")
		return
	}

	result.Success(ctx, "获取新书成功", books)
}

// GetBookDetail 获取图书详情
func (b *BookHandler) GetBookDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Log.Warn("GetBookDetail: 图书ID无效", zap.String("id", ctx.Param("id")), zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "图书ID无效")
		return
	}
	book, err := b.bookService.GetBookByID(ctx.Request.Context(), uint64(id))
	if err != nil {
		logger.Log.Error("GetBookDetail: 获取图书详情失败", zap.Int("id", id), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取图书详情失败")
		return
	}

	result.Success(ctx, "获取图书详情成功", book)
}

// SearchBooks
func (b *BookHandler) SearchBooks(ctx *gin.Context) {
	// 搜索栏输入文本
	keyword := ctx.Query("q")
	if keyword == "" {
		keyword = ctx.Query("keyword") // 兼容旧版本
	}
	if keyword == "" {
		logger.Log.Warn("SearchBooks: 搜索关键词为空")
		result.Fail(ctx, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	// 分页查询参数: page, page_size
	var pageReq request.BooksPageDTO
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		logger.Log.Warn("SearchBooks: 查询参数绑定失败，使用默认值", zap.Error(err))
		pageReq.Page = 1
		pageReq.PageSize = 12
	}
	// service 层返回分页结果
	pageResult, err := b.bookService.SearchBooksWithPagination(ctx.Request.Context(), keyword, &pageReq)
	if err != nil {
		logger.Log.Error("SearchBooks: 搜索图书失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "搜索图书失败")
		return
	}
	// 封装 VO
	result.Success(ctx, "搜索图书成功", &response.BooksPageVO{
		Books:     pageResult.Records,
		Total:     pageResult.Total,
		Page:      pageReq.Page,
		PageSize:  pageReq.PageSize,
		TotalPage: (pageResult.Total + int64(pageReq.PageSize) - 1) / int64(pageReq.PageSize),
	})
}

// GetBooksByCategory
func (b *BookHandler) GetBooksByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		logger.Log.Warn("GetBooksByCategory: 分类参数 category 为空")
		result.Fail(ctx, http.StatusBadRequest, "分类参数 category 不能为空")
		return
	}
	// 调用 service 层
	books, err := b.bookService.GetBooksByType(ctx.Request.Context(), category)
	if err != nil {
		logger.Log.Error("GetBooksByCategory: 获取分类图书失败", zap.String("category", category), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取分类图书失败")
		return
	}
	result.Success(ctx, "获取分类图书成功", books)
}
