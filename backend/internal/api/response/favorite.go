package response

import "github.com/wangn-tech/bookstore-go/internal/model"

type FavoritesPageVO struct {
	Favorites   []*model.Favorite `json:"favorites"`    // 一页显示的收藏列表
	Total       int64             `json:"total"`        // 总数量
	CurrentPage int               `json:"current_page"` // 当前页码
	TotalPages  int64             `json:"total_pages"`  // 总页数
	// PageSize    int               `json:"page_size"`    // 每页数量
}
