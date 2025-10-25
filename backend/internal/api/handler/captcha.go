package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/service"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
	"go.uber.org/zap"
)

// CaptchaHandler
type CaptchaHandler struct {
	captchaService service.ICaptchaService
}

// NewCaptchaHandler 构造函数
func NewCaptchaHandler(captchaService service.ICaptchaService) *CaptchaHandler {
	return &CaptchaHandler{
		captchaService: captchaService,
	}
}

// GenerateCaptcha 生成验证码
func (c *CaptchaHandler) GenerateCaptcha(ctx *gin.Context) {
	// 生成验证码
	captchaResp, err := c.captchaService.GenerateCaptcha(ctx.Request.Context())
	if err != nil {
		logger.Log.Warn("GenerateCaptcha: 生成验证码失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "生成验证码失败")
		return
	}

	result.Success(ctx, "生成验证码成功", captchaResp)
}

// VerifyCaptcha 验证验证码
func (c *CaptchaHandler) VerifyCaptcha(ctx *gin.Context) {
	// 绑定请求参数
	var req struct {
		CaptchaID    string `json:"captcha_id" binding:"required"`
		CaptchaValue string `json:"captcha_value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("VerifyCaptcha: 验证码参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 验证验证码
	isValid := c.captchaService.VerifyCaptcha(ctx.Request.Context(), req.CaptchaID, req.CaptchaValue)
	if !isValid {
		logger.Log.Warn("VerifyCaptcha: 验证码验证失败")
		result.Fail(ctx, http.StatusBadRequest, "验证码错误")
		return
	}

	result.Success(ctx, "验证码验证成功", nil)
}
