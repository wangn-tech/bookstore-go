package bookstore

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/api/handler"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/service"
)

type CarouselRouter struct {
	carouselService service.ICarouselService
}

func (c *CarouselRouter) InitCarouselRouter(router *gin.RouterGroup) {
	c.carouselService = service.NewCarouselService(
		repository.NewCarouselDao(database.DB),
	)
	carouselHandler := handler.NewCarouselHandler(c.carouselService)

	carouselGroup := router.Group("/carousel")
	{
		carouselGroup.GET("/list", carouselHandler.GetActiveCarousels)
		carouselGroup.GET("/:id", carouselHandler.GetCarouselByID)
		carouselGroup.POST("/create", carouselHandler.CreateCarousel)
		carouselGroup.PUT("/:id", carouselHandler.UpdateCarousel)
		carouselGroup.DELETE("/:id", carouselHandler.DeleteCarousel)
	}

}
