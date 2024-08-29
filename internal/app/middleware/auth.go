package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			uuidString := uuid.New().String()
			cookie, _ := crypt.EncodeValue(uuidString)

			c.SetCookie("Auth", cookie, 3600, "", "localhost", false, true)

			c.String(http.StatusUnauthorized, "")
			c.Abort()

			logger.Info("Uuid",
				zap.String("uuidString", uuidString),
				zap.String("cookie", cookie),
			)

			return
		}

		Uuid, _ := crypt.DecodeValue(auth)

		logger.Info("Auth user",
			zap.String("Uuid", Uuid),
			zap.String("auth.Value", auth),
		)

		c.Set("UserId", Uuid)

		c.Next()
	}
}
