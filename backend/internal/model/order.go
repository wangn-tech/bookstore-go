package model

import "time"

// Order 订单模型
type Order struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	UserID      uint64     `gorm:"not null" json:"user_id"`         // 用户ID
	OrderNo     string     `gorm:"not null;unique" json:"order_no"` // 订单号
	TotalAmount int        `gorm:"not null" json:"total_amount"`    // 订单总金额（分）
	Status      int        `gorm:"default:0" json:"status"`         // 订单状态：0-待支付，1-已支付，2-已取消
	IsPaid      bool       `gorm:"default:false" json:"is_paid"`    // 是否已支付
	PaymentTime *time.Time `json:"payment_time"`                    // 支付时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联字段
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

func (o *Order) TableName() string {
	return "orders"
}
