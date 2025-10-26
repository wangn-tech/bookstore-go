package model

import "time"

// Category 图书分类模型
type Category struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;unique" json:"name"`   // 分类名称
	Description string    `json:"description"`                   // 分类描述
	Icon        string    `json:"icon"`                          // 分类图标
	Color       string    `json:"color"`                         // 分类颜色
	Gradient    string    `json:"gradient"`                      // 渐变色彩
	Sort        int       `gorm:"default:0" json:"sort"`         // 排序权重
	IsActive    bool      `gorm:"default:true" json:"is_active"` // 是否启用
	BookCount   int       `gorm:"default:0" json:"book_count"`   // 该分类下的图书数量
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Category) TableName() string {
	return "categories"
}
