package http

import (
	"github.com/gin-gonic/gin"
	_ "songs-lib/docs"
	"songs-lib/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		services: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.POST("/", h.createSong)
			songs.GET("/", h.getAllSongs)
			songs.PUT("/:id", h.updateSong)
			songs.DELETE("/:id", h.deleteSong)
		}

		songDetail := api.Group("/song_detail")
		{
			songDetail.PUT("/:id", h.updateSongDetail)
		}
	}

	return router
}
