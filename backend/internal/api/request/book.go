package request

type BooksPageDTO struct {
	Page     int `form:"page" json:"page"`           // 当前页码
	PageSize int `form:"page_size" json:"page_size"` // 每页数量
}
