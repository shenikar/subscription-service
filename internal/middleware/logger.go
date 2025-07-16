package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenikar/subscription-service/internal/logger"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		log.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": latency,
			"client":  c.ClientIP(),
		}).Info("request completed")
	}
}
