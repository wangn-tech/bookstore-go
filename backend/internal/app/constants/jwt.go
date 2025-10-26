package constants

import "time"

const (
	Issuer       = "Bookstore_Issuer"
	AccessToken  = "access"
	RefreshToken = "refresh"
	// AccessTokenExpire 访问token过期时间
	AccessTokenExpire = 2 * time.Hour
	// RefreshTokenExpire 刷新token过期时间
	RefreshTokenExpire = 7 * 24 * time.Hour
)
