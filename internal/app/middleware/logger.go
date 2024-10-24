package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		defer logger.Sync()

		logger.Info("Request",
			zap.String("URL", c.Request.URL.String()),
			zap.String("Method", c.Request.Method),
		)

		c.Next()

		logger.Info("Response",
			zap.String("URL", c.Request.URL.String()),
			zap.Int("Code", c.Writer.Status()),
			zap.Int("Size", c.Writer.Size()),
			zap.Duration("Duration", time.Since(t)),
		)
	}
}
