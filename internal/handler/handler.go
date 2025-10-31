package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/handler/middleware"
	"github.com/olenka-91/DocsServer/internal/service"
	"github.com/sirupsen/logrus"
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

	//Кастомный валидатор пароля
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("password_complexity", entity.ValidatePasswordComplexity)
		if err != nil {
			logrus.Fatal("Failed to register validation:", err)
		}
	}

	// router.NoMethod(middleware.Wrap(h.handleMethodNotAllowed))
	router.Use(gin.Logger())

	router = h.InitExtraErrorHandlers(router)

	g := router.Group("/api")
	{
		g.POST("/auth", h.signIn)
		g.POST("/register", h.signUp)
		g.POST("/refresh", h.refreshToken)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		private.POST("/logout", h.logout)
	}

	private = router.Group("/api/docs")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("", h.getDocsList)
		private.HEAD("", h.getDocsList)
		private.GET("/:id", h.getDoc)
		private.HEAD("/:id", h.getDoc)
		private.POST("", h.postDoc)
		private.DELETE("/:id", h.deleteDoc)
	}

	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router

}

func (h *Handler) InitExtraErrorHandlers(r *gin.Engine) *gin.Engine {
	r.Use(gin.CustomRecovery(func(c *gin.Context, r any) {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: "Internal Server Error",
		})
	}))

	r.HandleMethodNotAllowed = true //откл стандартный обработчик

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, entity.ErrorResponse{
			Message: "Method Not Allowed",
		})
	})

	r.NoRoute((func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/*any") {
			c.JSON(http.StatusNotImplemented, entity.ErrorResponse{
				Message: "Not Implemented",
			})
		}

		c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: "Not Found",
		})
	}))
	return r
}
