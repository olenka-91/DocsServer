package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/olenka-91/DocsServer/internal/handler/middleware"
	"github.com/olenka-91/DocsServer/internal/service"
	//	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	//	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	//	_ "github.com/olenka-91/DocsServer/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{services: serv}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	g := router.Group("/api/docs")
	{
		g.GET("", middleware.Wrap(h.getDocsList))
		g.HEAD("", middleware.Wrap(h.getDocsList))
		g.GET("/:id", middleware.Wrap(h.getDoc))
		g.HEAD("/:id", middleware.Wrap(h.getDoc))
		//g.POST("", middleware.Wrap(h.getDocsList))
	}

	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router

}
