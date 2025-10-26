package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/api/request"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"gorm.io/gorm"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, dto *request.CreateOrderDTO) (*model.Order, error)
	GetOrderByID(ctx context.Context, id uint64) (*model.Order, error)
	GetUserOrders(ctx context.Context, userID uint64, dto *request.OrdersPageDTO) (*result.PageResult[*model.Order], error)
	PayOrder(ctx context.Context, id uint64) error
	GetOrderStatistics(ctx context.Context, userID uint64) (map[string]any, error)
}

type OrderServiceImpl struct {
	orderDao *repository.OrderDao
	bookDao  *repository.BookDao
}

func NewOrderService(orderDao *repository.OrderDao, bookDao *repository.BookDao) IOrderService {
	return &OrderServiceImpl{
		orderDao: orderDao,
		bookDao:  bookDao,
	}
}

// GetOrderStatistics 获取订单统计数据
func (o *OrderServiceImpl) GetOrderStatistics(ctx context.Context, userID uint64) (map[string]any, error) {
	return o.orderDao.GetOrderStatistics(ctx, userID)
}

// PayOrder 支付订单
func (o *OrderServiceImpl) PayOrder(ctx context.Context, id uint64) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 加锁读取订单及其明细，防并发
		order, err := o.orderDao.GetOrderWithItemsForUpdate(ctx, tx, id)
		if err != nil {
			return err
		}
		if order.IsPaid {
			return errors.New("订单已支付")
		}

		// 原子扣减每个商品库存并增加销量
		for _, item := range order.OrderItems {
			if err := o.bookDao.DecreaseStockAndIncreaseSaleTx(ctx, tx, item.BookID, item.Quantity); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("图书库存不足")
				}
				return err
			}
		}

		// 标记订单支付完成（幂等）
		if err := o.orderDao.MarkPaidTx(ctx, tx, id, time.Now()); err != nil {
			return err
		}
		return nil
	})
}

// GetUserOrders 获取用户的订单列表，支持分页
func (o *OrderServiceImpl) GetUserOrders(ctx context.Context, userID uint64, dto *request.OrdersPageDTO) (*result.PageResult[*model.Order], error) {
	return o.orderDao.GetUserOrdersByPage(ctx, userID, dto.Page, dto.PageSize)
}

// GetOrderByID 根据 ID 获取订单详情
func (o *OrderServiceImpl) GetOrderByID(ctx context.Context, id uint64) (*model.Order, error) {
	return o.orderDao.GetOrderByID(ctx, id)
}

// CreateOrder 创建订单
func (o *OrderServiceImpl) CreateOrder(ctx context.Context, dto *request.CreateOrderDTO) (*model.Order, error) {
	if len(dto.Items) == 0 {
		return nil, errors.New("订单项不能为空")
	}

	// 检查库存
	err := o.checkStockAvailability(ctx, dto.Items)
	if err != nil {
		return nil, err
	}
	// 生成订单号
	orderNo := o.generateOrderNo()
	// 计算总金额
	var totalAmount int               // 总金额
	var orderItems []*model.OrderItem // 订单项列表
	for _, item := range dto.Items {
		subtotal := item.Price * item.Quantity
		totalAmount += subtotal

		orderItems = append(orderItems, &model.OrderItem{
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: subtotal,
		})
	}
	// 创建订单
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      dto.UserID,
		TotalAmount: totalAmount,
		Status:      0,     // 待支付
		IsPaid:      false, // 未支付
	}
	// 使用事务创建订单和订单项
	err = o.orderDao.CreateOrderWithItems(ctx, order, orderItems)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// checkStockAvailability 检查订单项的库存是否充足
func (o *OrderServiceImpl) checkStockAvailability(ctx context.Context, items []request.CreateOrderItemDTO) error {
	for _, item := range items {
		book, err := o.bookDao.GetBookByIDForAdmin(ctx, item.BookID)
		if err != nil {
			return errors.New("图书不存在")
		}
		if book.Status != 1 {
			return errors.New("图书已下架")
		}
		if book.Stock < item.Quantity {
			return errors.New("库存不足")
		}
	}
	return nil
}

// generateOrderNo 生成唯一的订单号 (模拟实现)
func (o *OrderServiceImpl) generateOrderNo() string {
	orderNo := fmt.Sprintf("ORD%d", time.Now().UnixNano())
	return orderNo
}
