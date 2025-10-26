package bookstore

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/api/handler"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/middlerware"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/service"
)

type UserRouter struct {
	userService    service.IUserService
	captchaService service.ICaptchaService
}

func (u *UserRouter) InitUserRouter(router *gin.RouterGroup) {
	// 依赖注入
	u.userService = service.NewUserService(
		repository.NewUserDao(database.DB),
	)
	u.captchaService = service.NewCaptchaService()
	userHandler := handler.NewUserHandler(u.userService, u.captchaService)
	captchaHandler := handler.NewCaptchaHandler(u.captchaService)

	// "/user" 路由组
	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", userHandler.Register) // 用户注册
		userGroup.POST("/login", userHandler.Login)       // 用户登录

		// 需要认证的路由
		userGroup.Use(middlerware.JWTAuth())
		{
			userGroup.GET("/profile", userHandler.GetUserProfile)    // 获取用户信息
			userGroup.POST("/logout", userHandler.Logout)            // 用户登出
			userGroup.PUT("/profile", userHandler.UpdateUserProfile) // 更新用户信息
			userGroup.PUT("/password", userHandler.ChangePassword)   // 修改密码

		}
	}

	captchaGroup := router.Group("/captcha")
	{
		captchaGroup.GET("/generate", captchaHandler.GenerateCaptcha)
		captchaGroup.POST("/verify", captchaHandler.VerifyCaptcha)
	}
}
