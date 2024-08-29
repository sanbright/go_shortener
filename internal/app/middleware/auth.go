package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sanbright/go_shortener/internal/app/generator"
)

func Auth(crypt *generator.CryptGenerator, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := c.Cookie("Auth")
		if err != nil {
			logger.Info("Request Cookie Error",
				zap.String("ERROR", err.Error()),
			)

			for _, ck := range c.Request.Cookies() {
				logger.Info("Request Registered Cookie",
					zap.String("ck", ck.String()),
				)
			}

			c.String(http.StatusUnauthorized, "")
			c.Abort()
			return
		}

		UUID, _ := crypt.DecodeValue(auth)

		logger.Info("Auth user",
			zap.String("UUID", UUID),
			zap.String("auth.Value", auth),
		)
		c.Set("UserID", UUID)
		c.Next()
	}
}
