package repository

import (
	"context"
	"errors"

	"github.com/wangn-tech/bookstore-go/internal/model"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

// CheckUserExists 检查用户是否存在
func (u *UserDao) CheckUserExists(ctx context.Context, username, phone, email string) (bool, error) {
	// 只查询 id 字段
	q := u.db.WithContext(ctx).Model(&model.User{}).Select("id")

	// 动态构建查询条件
	applied := false
	if username != "" {
		q = q.Where("username = ?", username)
		applied = true
	}
	if phone != "" {
		if applied {
			q = q.Or("phone = ?", phone)
		} else {
			q = q.Where("phone = ?", phone)
			applied = true
		}
	}
	if email != "" {
		if applied {
			q = q.Or("email = ?", email)
		} else {
			q = q.Where("email = ?", email)
			applied = true
		}
	}

	if !applied {
		// 没有提供任何查询条件，返回 false
		return false, nil
	}

	// 只取一条即可判断存在
	var dst struct{ ID uint64 }
	err := q.Limit(1).Take(&dst).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateUser 创建新用户
func (u *UserDao) CreateUser(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

// GetUserByUsername 根据 username 获取 user
func (u *UserDao) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据 userID 获取 user
func (u *UserDao) GetUserByID(ctx context.Context, userID uint64) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
