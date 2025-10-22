package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"` // 是否为管理员
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
