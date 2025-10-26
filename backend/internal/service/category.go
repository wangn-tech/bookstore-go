package service

import (
	"context"

	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
)

type ICategoryService interface {
	GetAllCategories(ctx context.Context) ([]*model.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint64) error
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

// GetCategoryByID 根据 ID 获取分类详情
func (c *CategoryServiceImpl) GetCategoryByID(ctx context.Context, id uint64) (*model.Category, error) {
	return c.categoryDao.GetCategoryByID(ctx, id)
}

// CreateCategory 创建分类
func (c *CategoryServiceImpl) CreateCategory(ctx context.Context, category *model.Category) error {
	return c.categoryDao.CreateCategory(ctx, category)
}

// UpdateCategory 更新分类
func (c *CategoryServiceImpl) UpdateCategory(ctx context.Context, category *model.Category) error {
	return c.categoryDao.UpdateCategory(ctx, category)
}

// DeleteCategory 删除分类
func (c *CategoryServiceImpl) DeleteCategory(ctx context.Context, id uint64) error {
	return c.categoryDao.DeleteCategory(ctx, id)
}
