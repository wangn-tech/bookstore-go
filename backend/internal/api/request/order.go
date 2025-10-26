package request

// CreateOrderItemDTO 创建订单项请求
type CreateOrderItemDTO struct {
	BookID   uint64 `json:"book_id"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"` // 价格（元）
}

// CreateOrderDTO 创建订单请求
type CreateOrderDTO struct {
	UserID uint64               `json:"user_id" binding:"required"`
	Items  []CreateOrderItemDTO `json:"items" binding:"required"`
}

type OrdersPageDTO struct {
	Page     int `form:"page" json:"page"`           // 当前页码
	PageSize int `form:"page_size" json:"page_size"` // 每页数量
}
