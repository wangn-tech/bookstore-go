package repository

import (
	"context"

	"github.com/wangn-tech/bookstore-go/internal/model"
	"gorm.io/gorm"
)

type CategoryDao struct {
	db *gorm.DB
}

func NewCategoryDao(db *gorm.DB) *CategoryDao {
	return &CategoryDao{
		db: db,
	}
}

// GetAllCategories 获取所有分类
func (r *CategoryDao) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
