package handler

import (
	"github.com/gin-gonic/gin"
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

	router.GET("/api/docs", h.getDocsList)
	router.HEAD("/api/docs", h.getDocsList)

	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router

}
