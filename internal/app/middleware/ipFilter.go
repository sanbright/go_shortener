package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
)

// IpFilter middleware производит логгирование запроса и ответа
func IpFilter(logger *zap.Logger, trustedSubnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer logger.Sync()

		IP := c.GetHeader("X-Real-IP")
		if IP == "" {
			c.String(http.StatusForbidden, "")
			c.Abort()
			return
		}

		if !checkIP(IP, trustedSubnet) {
			c.String(http.StatusForbidden, "")
			c.Abort()
			return
		}

		c.Next()
	}
}

func checkIP(IP string, trustedSubnet string) bool {
	ip := net.ParseIP(IP)
	if ip == nil {
		return false
	}

	if _, ipNet, err := net.ParseCIDR(trustedSubnet); err == nil {
		if ipNet.Contains(ip) {
			return true
		}
	}

	ipList := strings.Split(trustedSubnet, ",")
	for _, validIP := range ipList {
		if _, ipNet, err := net.ParseCIDR(validIP); err == nil {
			if ipNet.Contains(ip) {
				return true
			}
		}
	}

	return false
}
