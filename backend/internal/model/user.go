package model

import (
	"time"

	"gorm.io/gorm"
)

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

// BeforeCreate 在创建记录前设置 CreatedAt 和 UpdatedAt 字段
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// BeforeUpdate 在更新记录前设置 UpdatedAt 字段
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
