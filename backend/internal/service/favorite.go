package service

import (
	"context"

	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/api/request"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
)

type IFavoriteService interface {
	AddFavorite(ctx context.Context, userID uint64, bookID uint64) error
	RemoveFavorite(ctx context.Context, userID uint64, bookID uint64) error
	IsFavorited(ctx context.Context, userID uint64, bookID uint64) (bool, error)
	GetUserFavorites(ctx context.Context, userID uint64, pageReq *request.FavoritePageDTO, timeFilter string) (*result.PageResult[*model.Favorite], error)
	GetUserFavoriteCount(ctx context.Context, userID uint64) (int64, error)
}

type FavoriteServiceImpl struct {
	favoriteDao repository.FavoriteDao
}

func NewFavoriteService(favoriteDao repository.FavoriteDao) IFavoriteService {
	return &FavoriteServiceImpl{
		favoriteDao: favoriteDao,
	}
}

// AddFavorite 添加收藏
func (s *FavoriteServiceImpl) AddFavorite(ctx context.Context, userID uint64, bookID uint64) error {
	return s.favoriteDao.AddFavorite(ctx, userID, bookID)
}

// RemoveFavorite 删除收藏
func (s *FavoriteServiceImpl) RemoveFavorite(ctx context.Context, userID uint64, bookID uint64) error {
	return s.favoriteDao.RemoveFavorite(ctx, userID, bookID)
}

// IsFavorited 检查是否已收藏
func (s *FavoriteServiceImpl) IsFavorited(ctx context.Context, userID uint64, bookID uint64) (bool, error) {
	return s.favoriteDao.IsFavorited(ctx, userID, bookID)
}

// GetUserFavorites 获取用户收藏列表，支持分页和时间过滤
func (s *FavoriteServiceImpl) GetUserFavorites(ctx context.Context, userID uint64, pageReq *request.FavoritePageDTO, timeFilter string) (*result.PageResult[*model.Favorite], error) {
	// 这里可以根据 timeFilter 参数添加时间过滤逻辑
	return s.favoriteDao.GetUserFavorites(ctx, userID, pageReq.Page, pageReq.PageSize, timeFilter)
}

// GetUserFavoriteCount 获取用户收藏数量
func (s *FavoriteServiceImpl) GetUserFavoriteCount(ctx context.Context, userID uint64) (int64, error) {
	return s.favoriteDao.GetUserFavoriteCount(ctx, userID)
}
