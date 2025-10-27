package repository

import (
	"context"

	"github.com/wangn-tech/bookstore-go/internal/model"
	"gorm.io/gorm"
)

type CarouselDao struct {
	db *gorm.DB
}

func NewCarouselDao(db *gorm.DB) *CarouselDao {
	return &CarouselDao{
		db: db,
	}
}

// DeleteCarouselByID 删除轮播图
func (dao *CarouselDao) DeleteCarouselByID(ctx context.Context, id uint64) error {
	return dao.db.WithContext(ctx).Delete(&model.Carousel{}, id).Error
}

// UpdateCarousel 更新轮播图
func (dao *CarouselDao) UpdateCarousel(ctx context.Context, carousel *model.Carousel) error {
	return dao.db.WithContext(ctx).Save(carousel).Error
}

// GetActiveCarousels 获取所有激活的轮播图
func (dao *CarouselDao) GetActiveCarousels(ctx context.Context) ([]*model.Carousel, error) {
	var carousels []*model.Carousel
	if err := dao.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&carousels).Error; err != nil {
		return nil, err
	}
	return carousels, nil
}

// GetCarouselByID 获取指定ID的轮播图
func (dao *CarouselDao) GetCarouselByID(ctx context.Context, id uint64) (*model.Carousel, error) {
	var carousel model.Carousel
	if err := dao.db.WithContext(ctx).First(&carousel, id).Error; err != nil {
		return nil, err
	}
	return &carousel, nil
}

// CreateCarousel 创建轮播图
func (dao *CarouselDao) CreateCarousel(ctx context.Context, carousel *model.Carousel) error {
	return dao.db.WithContext(ctx).Create(carousel).Error
}
