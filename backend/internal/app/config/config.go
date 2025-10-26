package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

var AppConf *Config

// Config 是应用程序的主配置结构体
type Config struct {
	Server ServerConfig   `mapstructure:"server"`
	MySQL  DatabaseConfig `mapstructure:"mysql"`
	Redis  RedisConfig    `mapstructure:"redis"`
	JWT    JWTConfig      `mapstructure:"jwt"`
	Log    LogConfig      `mapstructure:"log"`
}

// ServerConfig 后端服务端口配置
type ServerConfig struct {
	Port      int    `mapstructure:"port"`       // 服务端口
	AdminPort string `mapstructure:"admin_port"` // 管理端口
	Mode      string `mapstructure:"mode"`       // 运行模式: debug, release, test
}

// DatabaseConfig 数据库 MySQL 配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// LogConfig 定义了日志的配置参数
type LogConfig struct {
	Level      string `mapstructure:"level"`      // 日志级别, 例如: debug, info, warn, error
	Format     string `mapstructure:"format"`     // 日志格式, 例如: console, json
	ShowLine   bool   `mapstructure:"show-line"`  // 是否显示行号
	Stacktrace bool   `mapstructure:"stacktrace"` // 是否开启堆栈跟踪
}

// Init 读取解析配置文件
func Init() {
	// 解析命令行参数 “env”, 默认为 "dev"
	env := pflag.String("env", "dev", "Specify the environment config to use: [dev, prod, test]")
	pflag.Parse()

	// 设置配置文件: ./config/config-{env}.yaml
	config := viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName(fmt.Sprintf("config-%s", *env))
	config.SetConfigType("yaml")

	// 读取配置文件
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 解析配置
	AppConf = &Config{}
	if err := config.Unmarshal(AppConf); err != nil {
		panic(fmt.Errorf("unable to decode config: %w", err))
	}

	log.Println("Configuration loaded successfully")
}
