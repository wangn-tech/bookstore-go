package request

// UserRegisterDTO 注册请求结构体
type UserRegisterDTO struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone" binding:"required"`
	CaptchaID       string `json:"captcha_id" binding:"required"`
	CaptchaValue    string `json:"captcha_value" binding:"required"`
}

// UserLoginDTO 登录请求结构体
type UserLoginDTO struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	CaptchaID    string `json:"captcha_id" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}
