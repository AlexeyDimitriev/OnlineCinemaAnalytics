package http

import (
	"online-cinema-analytics/internal/infrastructure/http/handler"

	"github.com/gin-gonic/gin"
)

func Run(addr string, evtHandler *handler.EventHandler) error {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context)  {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/events", evtHandler.Create)
	}
	return r.Run(addr)
}
