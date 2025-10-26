package repository

import (
	"context"
	"time"

	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderDao struct {
	db *gorm.DB
}

func NewOrderDao(db *gorm.DB) *OrderDao {
	return &OrderDao{
		db: db,
	}
}

// GetOrderStatistics 获取订单统计数据
func (o *OrderDao) GetOrderStatistics(ctx context.Context, userID uint64) (map[string]any, error) {
	var stats struct {
		TotalOrders int64   `json:"total_orders"`
		TotalAmount float64 `json:"total_amount"` // 这里有 bug, 正常存储 RMB 应该用 int 存储分单位
		// TotalAmount   int     `json:"total_amount"`
		PaidOrders    int64 `json:"paid_orders"`
		PendingOrders int64 `json:"pending_orders"`
	}

	// 单次查询完成总数、总金额、已支付、待支付（避免 N+1/多次往返；COALESCE 处理无单时 SUM 为 NULL）
	err := o.db.WithContext(ctx).Model(&model.Order{}).
		Select(
			"COUNT(*) AS total_orders",
			"COALESCE(SUM(total_amount), 0) AS total_amount",
			"SUM(CASE WHEN is_paid = TRUE THEN 1 ELSE 0 END) AS paid_orders",
			"COUNT(*) - SUM(CASE WHEN is_paid = TRUE THEN 1 ELSE 0 END) AS pending_orders",
		).
		Where("user_id = ?", userID).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"total_orders":   stats.TotalOrders,
		"total_amount":   stats.TotalAmount,
		"paid_orders":    stats.PaidOrders,
		"pending_orders": stats.PendingOrders,
	}
	return result, nil

}

// GetOrderWithItemsForUpdate 在事务中查询订单并锁行（含订单项）
func (o *OrderDao) GetOrderWithItemsForUpdate(ctx context.Context, tx *gorm.DB, id uint64) (*model.Order, error) {
	db := tx
	if db == nil {
		db = o.db
	}
	var order model.Order
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}). // FOR UPDATE 行级锁
		Preload("OrderItems").                       // 预加载 OrderItems，不锁
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// MarkPaidTx 标记订单已支付（幂等：仅当未支付才更新）
func (o *OrderDao) MarkPaidTx(ctx context.Context, tx *gorm.DB, id uint64, paidAt time.Time) error {
	db := tx
	if db == nil {
		db = o.db
	}
	res := db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? AND is_paid = ?", id, false). // 仅当未支付时更新, is_paid: false -> true
		Updates(map[string]any{
			"status":       1,
			"is_paid":      true,
			"payment_time": paidAt,
		})
	if res.Error != nil {
		return res.Error
	}
	// RowsAffected 为 0 说明已支付或不存在
	return nil
}

// GetUserOrdersByPage 获取用户的订单列表，支持分页
func (o *OrderDao) GetUserOrdersByPage(ctx context.Context, userID uint64, page int, pageSize int) (*result.PageResult[*model.Order], error) {
	var total int64
	var orders []*model.Order

	// 构建查询
	query := o.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询，关联预加载 Preload
	// Preload 是 GORM 的“预加载关联”，用额外的查询把关联数据一次性取回并填充到结构体字段，避免 N+1 查询
	// Preload("OrderItems.Book") 表示两级预加载：先加载每个订单的 OrderItems，再把每个 OrderItem 关联的 Book 一并加载
	// 执行过程:
	// 	// 查询 orders: SELECT ... FROM orders WHERE user_id = ? LIMIT ? OFFSET ?
	// 	// 查询 order_items: SELECT ... FROM order_items WHERE order_id IN (订单 ID 集合)
	// 	// 查询 books: SELECT ... FROM books WHERE id IN (order_items 中的 book_id 集合)
	if err := query.Preload("OrderItems.Book").
		Scopes(result.Paginate(&page, &pageSize)).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return &result.PageResult[*model.Order]{
		Total:   total,
		Records: orders,
	}, nil
}

// GetOrderByID 根据 ID 获取订单详情
func (o *OrderDao) GetOrderByID(ctx context.Context, id uint64) (*model.Order, error) {
	var order model.Order
	// 关联预加载 Preload
	err := o.db.WithContext(ctx).Preload("OrderItems.Book").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateOrder 创建订单
func (o *OrderDao) CreateOrder(ctx context.Context, order *model.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

// CreateOrderWithItems 创建订单及其订单项，使用事务
func (o *OrderDao) CreateOrderWithItems(ctx context.Context, order *model.Order, items []*model.OrderItem) error {
	return o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// 创建订单项
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
