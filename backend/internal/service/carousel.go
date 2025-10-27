package service

import (
	"context"

	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
)

type ICarouselService interface {
	GetActiveCarousels(ctx context.Context) ([]*model.Carousel, error)
	GetCarouselByID(ctx context.Context, id uint64) (*model.Carousel, error)
	CreateCarousel(ctx context.Context, carousel *model.Carousel) error
	UpdateCarousel(ctx context.Context, carousel *model.Carousel) error
	DeleteCarousel(ctx context.Context, id uint64) error
}

type CarouselServiceImpl struct {
	carouselDao *repository.CarouselDao
}

func NewCarouselService(carouselDao *repository.CarouselDao) ICarouselService {
	return &CarouselServiceImpl{
		carouselDao: carouselDao,
	}
}

// CreateCarousel 创建轮播图
func (s *CarouselServiceImpl) CreateCarousel(ctx context.Context, carousel *model.Carousel) error {
	return s.carouselDao.CreateCarousel(ctx, carousel)
}

// GetCarouselByID 获取指定ID的轮播图
func (s *CarouselServiceImpl) GetCarouselByID(ctx context.Context, id uint64) (*model.Carousel, error) {
	return s.carouselDao.GetCarouselByID(ctx, id)
}

// GetActiveCarousels 获取所有激活的轮播图
func (s *CarouselServiceImpl) GetActiveCarousels(ctx context.Context) ([]*model.Carousel, error) {
	return s.carouselDao.GetActiveCarousels(ctx)
}

// UpdateCarousel 更新轮播图
func (s *CarouselServiceImpl) UpdateCarousel(ctx context.Context, carousel *model.Carousel) error {
	return s.carouselDao.UpdateCarousel(ctx, carousel)
}

// DeleteCarousel 删除轮播图
func (s *CarouselServiceImpl) DeleteCarousel(ctx context.Context, id uint64) error {
	return s.carouselDao.DeleteCarouselByID(ctx, id)

}
