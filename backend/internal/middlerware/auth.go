package middlerware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/app/constants"
	"github.com/wangn-tech/bookstore-go/internal/utils"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
	"go.uber.org/zap"
)

// JWTAuth 创建一个 JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取 token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warn("请求未携带 token，无权限访问", zap.String("path", ctx.Request.URL.Path))
			result.Fail(ctx, http.StatusUnauthorized, "请求未携带token，无权限访问")
			ctx.Abort()
			return
		}

		// 解析 token
		// 检查 "Bearer " 前缀
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Log.Warn("请求 token 格式错误", zap.String("path", ctx.Request.URL.Path))
			result.Fail(ctx, http.StatusUnauthorized, "请求 token 格式错误")
			ctx.Abort()
			return
		}
		tokenStr := tokenParts[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			// token 解析失败
			logger.Log.Warn("token 解析失败", zap.Error(err), zap.String("path", ctx.Request.URL.Path))
			// 根据不同的错误类型返回不同的消息
			result.Fail(ctx, http.StatusUnauthorized, "无效的token")
			ctx.Abort()
			return
		}
		valid := utils.IsTokenValidInRedis(claims.UserID, authHeader, constants.AccessToken)
		if !valid {
			logger.Log.Warn("JWTAuth: token 已失效或被撤销",
				zap.Uint64("userID", claims.UserID),
				zap.Error(err))
			result.Fail(ctx, http.StatusUnauthorized, "token已失效，请重新登录")
			ctx.Abort()
			return
		}
		// 将当前请求的 claims 信息保存到请求的上下文 c 上
		ctx.Set(constants.UserID, claims.UserID)
		ctx.Set(constants.Username, claims.Username)
		ctx.Next()
	}
}
