package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/api/request"
	"github.com/wangn-tech/bookstore-go/internal/api/response"
	"github.com/wangn-tech/bookstore-go/internal/service"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
	"go.uber.org/zap"
)

// UserHandler
type UserHandler struct {
	userService    service.IUserService
	captchaService service.ICaptchaService
}

// NewUserHandler 构造函数
func NewUserHandler(userService service.IUserService, captchaService service.ICaptchaService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		captchaService: captchaService,
	}
}

// Register 用户注册
func (u *UserHandler) Register(ctx *gin.Context) {
	var req request.UserRegisterDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Register: 用户注册参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}
	// 验证密码和确认密码是否匹配
	if req.Password != req.ConfirmPassword {
		logger.Log.Warn("Register: 密码和确认密码不匹配")
		result.Fail(ctx, http.StatusBadRequest, "密码和确认密码不匹配")
		return
	}

	// 验证验证码
	if !u.captchaService.VerifyCaptcha(ctx.Request.Context(), req.CaptchaID, req.CaptchaValue) {
		logger.Log.Warn("Register: 验证码验证失败")
		result.Fail(ctx, http.StatusBadRequest, "验证码错误")
		return
	}

	// 调用 service 层注册用户
	err := u.userService.Register(ctx.Request.Context(), req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		logger.Log.Warn("Register: 用户注册失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "用户注册失败")
		return
	}

	// 注册成功
	result.Success(ctx, "用户注册成功", nil)
}

// Login 用户登录
func (u *UserHandler) Login(ctx *gin.Context) {
	var req request.UserLoginDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Login: 用户登录参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 验证验证码
	if !u.captchaService.VerifyCaptcha(ctx.Request.Context(), req.CaptchaID, req.CaptchaValue) {
		logger.Log.Warn("Login: 验证码验证失败")
		result.Fail(ctx, http.StatusBadRequest, "验证码错误")
		return
	}

	// 调用 service 层登录用户
	userLoginVO, err := u.userService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		logger.Log.Warn("Login: 用户登录失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "用户登录失败")
		return
	}

	// 登录成功
	result.Success(ctx, "用户登录成功", userLoginVO)
}

// GetUserProfile 获取用户信息
func (u *UserHandler) GetUserProfile(ctx *gin.Context) {
	// 从 context 中获取 userID
	userIDVal, ok := ctx.Get("userID")
	if !ok {
		logger.Log.Warn("GetUserProfile: 用户ID不存在")
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		logger.Log.Warn("GetUserProfile: 用户ID类型不匹配", zap.Any("userID", userIDVal))
		result.Fail(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	// 调用 service 层获取用户信息
	user, err := u.userService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		logger.Log.Warn("GetUserProfile: 获取用户信息失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 返回用户信息
	result.Success(ctx, "获取用户信息成功", &response.UserProfileVO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}
