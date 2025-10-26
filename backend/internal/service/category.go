package service

import (
	"context"

	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
)

type ICategoryService interface {
	GetAllCategories(ctx context.Context) ([]*model.Category, error)
}

type CategoryServiceImpl struct {
	// Define fields, e.g., repository interfaces
	categoryDao *repository.CategoryDao
}

func NewCategoryService(categoryDao *repository.CategoryDao) ICategoryService {
	return &CategoryServiceImpl{
		categoryDao: categoryDao,
	}

}

// GetAllCategories 获取所有分类
func (c *CategoryServiceImpl) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	return c.categoryDao.GetAllCategories(ctx)
}
