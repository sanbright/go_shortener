package middleware

import (
	"sanbright/go_shortener/internal/app/generator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthGen middleware проверяет по Cookie авторизован ли пользователь, если нет то авторизует его
func AuthGen(crypt *generator.CryptGenerator, domain string, logger *zap.Logger) gin.HandlerFunc {
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

			uuidString := uuid.New().String()
			auth, _ = crypt.EncodeValue(uuidString)
			c.SetCookie("Auth", auth, 200000, "", domain, false, true)

			logger.Info("Set Cookie",
				zap.String("uuidString", uuidString),
				zap.String("cookie", auth),
			)
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
