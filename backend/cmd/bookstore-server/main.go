package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/app/config"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/redis"
	"github.com/wangn-tech/bookstore-go/internal/router"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
)

func main() {
	// 配置文件初始化 ./config/config-dev.yaml
	config.Init()
	// 初始化 Logger
	logger.InitLogger()

	// 初始化数据库
	database.InitDB()
	// 初始化 Redis
	redis.InitRedis()

	// 初始化 *gin.Engine
	gin.SetMode(config.AppConf.Server.Mode)
	r := gin.Default()
	// 注册路由
	router.InitRouter(r)

	// 启动服务
	port := fmt.Sprintf(":%d", config.AppConf.Server.Port)
	logger.Log.Info("Server starting on port " + port)
	if err := r.Run(port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
