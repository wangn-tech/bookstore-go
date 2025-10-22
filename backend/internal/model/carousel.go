package model

import (
	"time"
)

// Carousel 轮播图模型
type Carousel struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null;comment:轮播图标题"`
	Description string    `json:"description" gorm:"type:text;comment:轮播图描述"`
	ImageURL    string    `json:"image_url" gorm:"not null;comment:轮播图图片URL"`
	LinkURL     string    `json:"link_url" gorm:"comment:点击跳转链接"`
	SortOrder   int       `json:"sort_order" gorm:"default:0;comment:排序"`
	IsActive    bool      `json:"is_active" gorm:"default:true;comment:是否激活"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Carousel) TableName() string {
	return "carousel"
}
