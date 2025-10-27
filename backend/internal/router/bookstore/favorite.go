package bookstore

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/bookstore-go/internal/api/handler"
	"github.com/wangn-tech/bookstore-go/internal/app/initializer/database"
	"github.com/wangn-tech/bookstore-go/internal/middlerware"
	"github.com/wangn-tech/bookstore-go/internal/repository"
	"github.com/wangn-tech/bookstore-go/internal/service"
)

type FavoriteRouter struct {
	favoriteService service.IFavoriteService
}

func (f *FavoriteRouter) InitFavoriteRouter(router *gin.RouterGroup) {
	f.favoriteService = service.NewFavoriteService(
		repository.NewFavoriteDao(database.DB),
	)
	favoriteHandler := handler.NewFavoriteHandler(f.favoriteService)

	favoriteGroup := router.Group("/favorite")
	favoriteGroup.Use(middlerware.JWTAuth())
	{
		favoriteGroup.POST("/:id", favoriteHandler.AddFavorite)           // Add a book to favorites
		favoriteGroup.DELETE("/:id", favoriteHandler.RemoveFavorite)      // Remove a book from favorites
		favoriteGroup.GET("/:id/check", favoriteHandler.CheckFavorite)    // Check if a book is in favorites
		favoriteGroup.GET("/list", favoriteHandler.GetUserFavorites)      // Get all favorite books for a user
		favoriteGroup.GET("/count", favoriteHandler.GetUserFavoriteCount) // Get the count of favorite books for a user
	}
}
