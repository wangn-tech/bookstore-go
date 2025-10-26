package result

import (
	"gorm.io/gorm"
)

const (
	// MinPageSize 最小单页数量
	MinPageSize = 12
	// MaxPageSize 最大单页数量
	MaxPageSize = 100
)

// PageResult 分页结果
type PageResult[T any] struct {
	Total   int64 `json:"total"`   //总记录数
	Records []T   `json:"records"` //当前页数据集合
}

// PageVerify 分页参数校验, 修正非法参数
func PageVerify(page *int, pageSize *int) {
	// 过滤 当前页、单页数量
	if *page < 1 {
		*page = 1
	}
	switch {
	case *pageSize > 100:
		*pageSize = MaxPageSize
	case *pageSize <= 0:
		*pageSize = MinPageSize
	}
}

// Paginate
//
//	// 分页参数校验
//	// 拼接分页 sql
func Paginate(page *int, pageSize *int) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		// 分页校验
		PageVerify(page, pageSize)

		// 拼接分页
		d.Offset((*page - 1) * *pageSize).Limit(*pageSize)
		return d
	}
}
