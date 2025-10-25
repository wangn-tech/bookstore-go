package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/api/request"
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
