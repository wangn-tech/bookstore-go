package service

import "context"

type IUserService interface {
	// Login 登录
	Login(ctx context.Context) error
	// Logout 登出
	Logout(ctx context.Context) error
	// Register 注册
	Register(ctx context.Context) error
	// GetUserProfile 获取用户信息
	GetUserProfile(ctx context.Context) error
	// UpdateUserProfile 更新用户信息
	UpdateUserProfile(ctx context.Context) error
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context) error
}
