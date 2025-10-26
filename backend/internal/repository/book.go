package repository

import (
	"context"

	"github.com/wangn-tech/bookstore-go/common/result"
	"github.com/wangn-tech/bookstore-go/internal/model"
	"gorm.io/gorm"
)

type BookDao struct {
	db *gorm.DB
}

func NewBookDao(db *gorm.DB) *BookDao {
	return &BookDao{
		db: db,
	}
}

// DecreaseStockAndIncreaseSaleTx 原子扣减库存并增加销量
func (b *BookDao) DecreaseStockAndIncreaseSaleTx(ctx context.Context, tx *gorm.DB, bookID uint64, qty int) error {
	db := tx
	if db == nil {
		db = b.db
	}
	// UPDATE books SET stock = stock - :qty, sale  = sale  + :qty WHERE id = :book_id AND stock >= :qty;
	// 单条 SQL 语句实现原子操作，避免超卖
	// 	// 两个并发请求同时扣同一本书的库存时，InnoDB 对该行加锁，先提交者成功更新并减少库存
	// 	// 后提交者因条件 stock >= qty 不再满足而更新 0 行，从而避免超卖
	res := db.WithContext(ctx).
		Model(&model.Book{}).
		Where("id = ? AND stock >= ?", bookID, qty).
		Updates(map[string]any{
			"stock": gorm.Expr("stock - ?", qty),
			"sale":  gorm.Expr("sale + ?", qty),
		})
	if res.Error != nil {
		return res.Error
	}
	// 没有受影响行表示库存不足或书不存在
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetBookByIDForAdmin 根据ID获取书籍信息（管理员用）, 不过滤 Status
func (b *BookDao) GetBookByIDForAdmin(ctx context.Context, id uint64) (*model.Book, error) {
	var book model.Book
	err := b.db.WithContext(ctx).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// GetBooksByPage 分页获取书籍列表
func (b *BookDao) GetBooksByPage(ctx context.Context, page int, pageSize int) (*result.PageResult[*model.Book], error) {
	var total int64
	var books []*model.Book

	// 构建查询
	query := b.db.WithContext(ctx).Model(&model.Book{})

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := query.Scopes(result.Paginate(&page, &pageSize)).Find(&books).Error; err != nil {
		return nil, err
	}

	return &result.PageResult[*model.Book]{
		Total:   total,
		Records: books,
	}, nil
}

// GetHotBooks 获取热销图书
func (b *BookDao) GetHotBooks(ctx context.Context, limit int) ([]*model.Book, error) {
	var books []*model.Book
	err := b.db.WithContext(ctx).Model(&model.Book{}).
		Where("status = ?", 1).
		Order("sale DESC").
		Limit(limit).
		Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

// GetNewBooks 获取新书
func (b *BookDao) GetNewBooks(ctx context.Context, limit int) ([]*model.Book, error) {
	var books []*model.Book
	err := b.db.WithContext(ctx).Model(&model.Book{}).
		Where("status = ?", 1).
		Order("created_at DESC").
		Limit(limit).
		Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

// GetBookByID 根据ID获取书籍详情
func (b *BookDao) GetBookByID(ctx context.Context, id uint64) (*model.Book, error) {
	var book model.Book
	err := b.db.WithContext(ctx).Where("status = ?", 1).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// SearchBooksWithPagination 根据关键词搜索书籍，支持分页
func (b *BookDao) SearchBooksWithPagination(ctx context.Context, keyword string, page int, pageSize int) (*result.PageResult[*model.Book], error) {
	var total int64
	var books []*model.Book

	// 构建查询
	query := b.db.WithContext(ctx).Model(&model.Book{}).
		Where("title LIKE ? OR author LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := query.Scopes(result.Paginate(&page, &pageSize)).Find(&books).Error; err != nil {
		return nil, err
	}

	return &result.PageResult[*model.Book]{
		Total:   total,
		Records: books,
	}, nil
}

// GetBooksByCategory 根据分类获取书籍列表
func (b *BookDao) GetBooksByType(ctx context.Context, bookType string) ([]*model.Book, error) {
	var books []*model.Book
	err := b.db.WithContext(ctx).Model(&model.Book{}).
		Where("status = ? AND type = ?", 1, bookType).
		Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}
