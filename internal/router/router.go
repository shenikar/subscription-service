package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shenikar/subscription-service/internal/handler"
	"github.com/shenikar/subscription-service/internal/middleware"
)

func SetupRouter(h *handler.SubscriptionHandler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	api := r.Group("/api/v1")
	{
		sub := api.Group("/subscriptions")
		{
			sub.POST("/", h.Create)
			sub.GET("/", h.GetAll)
			sub.GET("/:id", h.GetByID)
			sub.PUT("/:id", h.Update)
			sub.DELETE("/:id", h.Delete)
			sub.GET("/total", h.TotalPrice)
		}
	}

	return r
}
