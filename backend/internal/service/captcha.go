package service

import (
	"context"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/redis"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
	"go.uber.org/zap"
)

// ICaptchaService 验证码服务接口
type ICaptchaService interface {
	// GenerateCaptcha 生成验证码
	GenerateCaptcha(ctx context.Context) (*CaptchaResponse, error)
	// VerifyCaptcha 验证验证码
	VerifyCaptcha(ctx context.Context, captchaID, captchaValue string) bool
}

// CaptchaServiceImpl 验证码服务实现
type CaptchaServiceImpl struct {
	store base64Captcha.Store
}

// NewCaptchaService 构造函数
func NewCaptchaService() ICaptchaService {
	return &CaptchaServiceImpl{
		store: base64Captcha.DefaultMemStore,
	}
}

// CaptchaResponse 验证码响应结构体
type CaptchaResponse struct {
	CaptchaID     string `json:"captcha_id"`
	CaptchaBase64 string `json:"captcha_base64"`
}

const (
	// Redis 键前缀
	captchaRedisKeyPrefix = "captcha:"
	// Redis 写入/读取超时
	captchaRedisOpTimeout = 2 * time.Second
	// 验证码有效期
	captchaTTL = 2 * time.Minute
)

// GenerateCaptcha 生成验证码
func (c *CaptchaServiceImpl) GenerateCaptcha(ctx context.Context) (*CaptchaResponse, error) {
	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(
		80,  // height
		240, // width
		4,   // length
		0.7, // maxSkew 干扰强度
		80,  // dotCount 干扰数量
	)

	// 生成验证码
	captcha := base64Captcha.NewCaptcha(driver, c.store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, err
	}

	// 限定 Redis 写入超时
	tctx, cancel := context.WithTimeout(ctx, captchaRedisOpTimeout)
	defer cancel()
	// 将验证码存储到 Redis, 设置 2 min 过期
	redisKey := captchaRedisKeyPrefix + id
	err = redis.RedisClient.Set(tctx, redisKey, answer, captchaTTL).Err()
	if err != nil {
		return nil, err
	}

	return &CaptchaResponse{
		CaptchaID:     id,
		CaptchaBase64: b64s,
	}, nil
}

// VerifyCaptcha 验证验证码
func (c *CaptchaServiceImpl) VerifyCaptcha(ctx context.Context, captchaID, captchaValue string) bool {
	if captchaID == "" || captchaValue == "" {
		return false
	}

	// 限定 Redis 写入超时
	tctx, cancel := context.WithTimeout(ctx, captchaRedisOpTimeout)
	defer cancel()
	// 从 Redis 获取验证码 val
	redisKey := captchaRedisKeyPrefix + captchaID
	val, err := redis.RedisClient.Get(tctx, redisKey).Result()
	if err != nil {
		return false
	}
	// 比较用户输入的 captchaValue 和存储的 val 是否相等，验证后删除该验证码
	isValid := val == captchaValue
	// 验证成功: 从 Redis 中删除该验证码
	if isValid {
		if err := redis.RedisClient.Del(tctx, redisKey).Err(); err != nil {
			// 验证码删除失败，记录日志
			logger.Log.Warn("captcha delete failed", zap.String("key", redisKey), zap.Error(err))
		}
	}
	return isValid
}
