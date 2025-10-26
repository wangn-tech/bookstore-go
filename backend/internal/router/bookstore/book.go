package bookstore

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/api/handler"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/service"
)

type BookRouter struct {
	bookService service.IBookService
}

func (b *BookRouter) InitBookRouter(router *gin.RouterGroup) {
	b.bookService = service.NewBookService(
		repository.NewBookDao(database.DB),
	)
	bookHandler := handler.NewBookHandler(b.bookService)

	bookRouter := router.Group("/book")
	{
		bookRouter.GET("/list", bookHandler.GetBookList)                      // 获取图书列表
		bookRouter.GET("/hot", bookHandler.GetHotBooks)                       // 获取热销图书
		bookRouter.GET("/new", bookHandler.GetNewBooks)                       // 获取新书
		bookRouter.GET("/detail/:id", bookHandler.GetBookDetail)              // 获取图书详情
		bookRouter.GET("/search", bookHandler.SearchBooks)                    // 搜索图书
		bookRouter.GET("/category/:category", bookHandler.GetBooksByCategory) // 获取分类下的图书
	}

}
