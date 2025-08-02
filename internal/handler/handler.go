package handler

import (
	"net/http"
	"strings"

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

	// router.NoMethod(middleware.Wrap(h.handleMethodNotAllowed))
	router.Use(gin.Logger())

	router = h.InitExtraErrorHandlers(router)

	g := router.Group("/api/docs", h.userIdentity)
	{
		g.GET("", middleware.Wrap(h.getDocsList))
		g.HEAD("", middleware.Wrap(h.getDocsList))
		g.GET("/:id", middleware.Wrap(h.getDoc))
		g.HEAD("/:id", middleware.Wrap(h.getDoc))
		g.POST("", middleware.Wrap(h.postDoc))
		g.DELETE("/:id", middleware.Wrap(h.deleteDoc))
	}

	g = router.Group("/api")
	{
		g.POST("/auth", middleware.Wrap(h.autorizeUser))
		g.POST("/register", middleware.Wrap(h.registerUser))
		//g.DELETE("/auth/:id", middleware.Wrap(h.deleteDoc))
	}

	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router

}

func (h *Handler) InitExtraErrorHandlers(r *gin.Engine) *gin.Engine {
	r.Use(gin.CustomRecovery(func(c *gin.Context, r any) {
		middleware.Fail(c, http.StatusInternalServerError,
			http.StatusInternalServerError, "internal server error")
	}))

	r.HandleMethodNotAllowed = true //откл стандартный обработчик

	r.NoMethod(middleware.Wrap(func(c *gin.Context) (any, any, *middleware.ErrorResponse, int) {
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusMethodNotAllowed, Text: "method not allowed"}, http.StatusMethodNotAllowed
	}))

	r.NoRoute(middleware.Wrap(func(c *gin.Context) (any, any, *middleware.ErrorResponse, int) {
		// Пример того что пока не реализовано
		if strings.HasPrefix(c.Request.URL.Path, "/api/docs/v2") {
			return nil, nil, &middleware.ErrorResponse{Code: http.StatusNotImplemented, Text: "status not implemented"},
				http.StatusNotImplemented
		}

		return nil, nil, &middleware.ErrorResponse{Code: http.StatusNotFound, Text: "not found"},
			http.StatusNotFound
	}))
	return r
}
