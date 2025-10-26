package model

import "time"

type Book struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Author      string    `json:"author"`
	Price       int       `json:"price"`       // 价格（元）
	Discount    int       `json:"discount"`    // 折扣（百分比，100表示无折扣）
	Type        string    `json:"type"`        // 图书类型
	Stock       int       `json:"stock"`       // 库存数量
	Status      int       `json:"status"`      // 图书状态：0-下架，1-上架
	Description string    `json:"description"` // 图书描述
	CoverURL    string    `json:"cover_url"`
	ISBN        string    `json:"isbn"`         // ISBN号
	Publisher   string    `json:"publisher"`    // 出版社
	PublishDate string    `json:"publish_date"` // 出版日期
	Pages       int       `json:"pages"`        // 页数
	Language    string    `json:"language"`     // 语言
	Format      string    `json:"format"`       // 装帧格式
	CategoryID  uint64    `json:"category_id"`  // 分类ID
	Sale        int       `json:"sale"`         // 销售量
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Book) TableName() string {
	return "books"
}
