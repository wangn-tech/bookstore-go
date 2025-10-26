package model

import "time"

// OrderItem 订单项模型
type OrderItem struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	OrderID   uint64    `gorm:"not null" json:"order_id"` // 订单ID
	BookID    uint64    `gorm:"not null" json:"book_id"`  // 图书ID
	Quantity  int       `gorm:"not null" json:"quantity"` // 数量
	Price     int       `gorm:"not null" json:"price"`    // 单价（分）
	Subtotal  int       `gorm:"not null" json:"subtotal"` // 小计（分）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联字段
	Book *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}
