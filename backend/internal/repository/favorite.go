package repository

import (
	"context"

	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"gorm.io/gorm"
)

type FavoriteDao struct {
	db *gorm.DB
}

func NewFavoriteDao(db *gorm.DB) FavoriteDao {
	return FavoriteDao{
		db: db,
	}
}

// AddFavorite 添加收藏记录
func (f *FavoriteDao) AddFavorite(ctx context.Context, userID uint64, bookID uint64) error {
	favorite := &model.Favorite{
		UserID: userID,
		BookID: bookID,
	}
	return f.db.WithContext(ctx).Create(favorite).Error
}

// RemoveFavorite 删除收藏记录
func (f *FavoriteDao) RemoveFavorite(ctx context.Context, userID uint64, bookID uint64) error {
	err := f.db.WithContext(ctx).Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&model.Favorite{}).Error
	return err
}

// IsFavorited 检查是否已收藏
func (f *FavoriteDao) IsFavorited(ctx context.Context, userID uint64, bookID uint64) (bool, error) {
	var count int64
	err := f.db.WithContext(ctx).Model(&model.Favorite{}).
		Where("user_id = ? AND book_id = ?", userID, bookID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserFavorites 获取用户收藏列表，支持分页 (未完成时间过滤)
func (f *FavoriteDao) GetUserFavorites(ctx context.Context, userID uint64, page int, pageSize int, timeFilter string) (*result.PageResult[*model.Favorite], error) {
	var favorites []*model.Favorite
	query := f.db.WithContext(ctx).Model(&model.Favorite{}).Where("user_id = ?", userID)

	// 根据 timeFilter 添加时间过滤逻辑
	// 例如，timeFilter 可以是 "last_week", "last_month", "all" 等
	// 这里仅作为示例，实际实现可能需要根据具体需求调整
	// switch timeFilter {
	// case "last_week":
	// 	query = query.Where("created_at >= NOW() - INTERVAL 7 DAY")
	// case "last_month":
	// 	query = query.Where("created_at >= NOW() - INTERVAL 1 MONTH")
	// case "all":
	// 	// 不添加时间过滤
	// default:
	// 	// 如果传入了未知的 timeFilter，可以选择返回错误或忽略
	// }

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	// err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&favorites).Error
	err = query.Scopes(result.Paginate(&page, &pageSize)).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	pageResult := &result.PageResult[*model.Favorite]{
		Total:   total,
		Records: favorites,
	}
	return pageResult, nil
}

// GetUserFavoriteCount 获取用户收藏数量
func (f *FavoriteDao) GetUserFavoriteCount(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := f.db.WithContext(ctx).Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
