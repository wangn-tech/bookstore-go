package response

import "github.com/wangn-tech/bookstore-go/internal/model"

type BooksPageVO struct {
	Books     []*model.Book `json:"books"`      // 一页显示的书籍列表
	Total     int64         `json:"total"`      // 总数量
	Page      int           `json:"page"`       // 当前页码
	PageSize  int           `json:"page_size"`  // 每页数量
	TotalPage int64         `json:"total_page"` // 总页数
}
