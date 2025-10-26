package service

import (
	"context"

	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/api/request"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
)

type IBookService interface {
	GetBooksByPage(ctx context.Context, dto request.BooksPageDTO) (*result.PageResult[*model.Book], error)
	GetBookByID(ctx context.Context, id uint64) (*model.Book, error)
	SearchBooksWithPagination(ctx context.Context, keyword string, pageReq *request.BooksPageDTO) (*result.PageResult[*model.Book], error)
	GetBooksByType(ctx context.Context, bookType string) ([]*model.Book, error)
	GetHotBooks(ctx context.Context, limit int) ([]*model.Book, error)
	GetNewBooks(ctx context.Context, limit int) ([]*model.Book, error)
}

type BookServiceImpl struct {
	bookDao *repository.BookDao
}

func NewBookService(bookDao *repository.BookDao) IBookService {
	return &BookServiceImpl{
		bookDao: bookDao,
	}
}

// GetBooksByPage 分页获取书籍列表
func (b *BookServiceImpl) GetBooksByPage(ctx context.Context, dto request.BooksPageDTO) (*result.PageResult[*model.Book], error) {
	pageResult, err := b.bookDao.GetBooksByPage(ctx, dto.Page, dto.PageSize)
	if err != nil {
		return nil, err
	}
	return pageResult, nil
}

// GetHotBooks 获取热销图书
func (b *BookServiceImpl) GetHotBooks(ctx context.Context, limit int) ([]*model.Book, error) {
	return b.bookDao.GetHotBooks(ctx, limit)
}

// GetNewBooks 获取新书
func (b *BookServiceImpl) GetNewBooks(ctx context.Context, limit int) ([]*model.Book, error) {
	return b.bookDao.GetNewBooks(ctx, limit)
}

// GetBookByID 根据ID获取书籍详情
func (b *BookServiceImpl) GetBookByID(ctx context.Context, id uint64) (*model.Book, error) {
	return b.bookDao.GetBookByID(ctx, id)
}

// SearchBooksWithPagination 根据关键词搜索书籍，支持分页
func (b *BookServiceImpl) SearchBooksWithPagination(ctx context.Context, keyword string, pageReq *request.BooksPageDTO) (*result.PageResult[*model.Book], error) {
	return b.bookDao.SearchBooksWithPagination(ctx, keyword, pageReq.Page, pageReq.PageSize)
}

// GetBooksByCategory 根据分类获取书籍列表
func (b *BookServiceImpl) GetBooksByType(ctx context.Context, bookType string) ([]*model.Book, error) {
	return b.bookDao.GetBooksByType(ctx, bookType)
}
