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
		categoryGroup.GET("/list", categoryHandler.GetCategories)
		categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
		categoryGroup.POST("/create", categoryHandler.CreateCategory)
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}

}
