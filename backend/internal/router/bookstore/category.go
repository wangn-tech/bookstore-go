package bookstore

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/api/handler"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/service"
)

type CategoryRouter struct {
	categoryService service.ICategoryService
}

func (c *CategoryRouter) InitCategoryRouter(router *gin.RouterGroup) {
	c.categoryService = service.NewCategoryService(
		repository.NewCategoryDao(database.DB),
	)
	categoryHandler := handler.NewCategoryHandler(c.categoryService)

	categoryGroup := router.Group("/category")
	{
		categoryGroup.GET("/list", categoryHandler.GetCategories)     // 获取所有分类
		categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)    // 根据 ID 获取分类详情
		categoryGroup.POST("/create", categoryHandler.CreateCategory) // 创建分类
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)     // 更新分类
		categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)  // 删除分类
	}

}
