package service

import (
	"context"
	"errors"

	"github.com/wangn-tech/bookstore-go/internal/api/response"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/utils"
)

type IUserService interface {
	// Login 登录
	Login(ctx context.Context, username, password string) (*response.UserLoginVO, error)
	// Logout 登出
	Logout(ctx context.Context, userID uint64) error
	// Register 注册
	Register(ctx context.Context, username, password, email, phone string) error
	// GetUserByID 获取用户信息
	GetUserByID(ctx context.Context, userID uint64) (*model.User, error)
	// UpdateUserProfile 更新用户信息
	UpdateUserProfile(ctx context.Context) error
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context) error
}

type UserServiceImpl struct {
	userDao *repository.UserDao
}

func NewUserService(userDao *repository.UserDao) IUserService {
	return &UserServiceImpl{
		userDao: userDao,
	}
}

// Login 用户登录
func (u *UserServiceImpl) Login(ctx context.Context, username, password string) (*response.UserLoginVO, error) {
	// 获取 user 信息
	user, err := u.userDao.GetUserByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}
	// 校验密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成 JWT Token 对
	tokenResponse, err := utils.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成 token 失败")
	}

	// 返回登录信息
	return &response.UserLoginVO{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
		UserInfo: &response.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}, nil
}

func (u *UserServiceImpl) Logout(ctx context.Context, userID uint64) error {
	// 调用 RevokeToken 函数，从 Redis 中删除用户的 token
	if err := utils.RevokeToken(userID); err != nil {
		return err
	}
	return nil
}

// Register 用户注册
func (u *UserServiceImpl) Register(ctx context.Context, username, password, email, phone string) error {
	// 查询用户是否已存在 (username, phone, email)
	exists, err := u.checkUserExists(ctx, username, phone, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名、手机号或邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户
	err = u.createUser(ctx, username, hashedPassword, email, phone)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) GetUserByID(ctx context.Context, userID uint64) (*model.User, error) {
	user, err := u.userDao.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}
	return user, nil
}

func (u *UserServiceImpl) UpdateUserProfile(ctx context.Context) error {
	return nil
}

func (u *UserServiceImpl) ChangePassword(ctx context.Context) error {
	return nil
}

// checkUserExists 检查用户是否已存在
func (u *UserServiceImpl) checkUserExists(ctx context.Context, username, phone, email string) (bool, error) {
	return u.userDao.CheckUserExists(ctx, username, phone, email)
}

// createUser 创建新用户
func (u *UserServiceImpl) createUser(ctx context.Context, username, hashedPassword, email, phone string) error {
	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Phone:    phone,
	}
	return u.userDao.CreateUser(ctx, user)
}
