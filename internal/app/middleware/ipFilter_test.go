package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}

func TestIpFilter(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		subnets string
		want    int
	}{
		{
			name:    "IP_not_exists",
			ip:      "127.0.32.1",
			subnets: "127.0.0.1/24",
			want:    http.StatusForbidden,
		},
		{
			name:    "IP_exists",
			ip:      "127.0.0.5",
			subnets: "127.0.0.1/24",
			want:    http.StatusOK,
		},
		{
			name:    "IP_empty",
			ip:      "",
			subnets: "127.0.0.1/24",
			want:    http.StatusForbidden,
		},
		{
			name:    "IP_many_subnet",
			ip:      "127.0.30.33",
			subnets: "127.0.0.1/24,127.0.30.1/24",
			want:    http.StatusOK,
		},
	}

	log := setupLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := IpFilter(log, tt.subnets)

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			headers := http.Header{}
			headers.Set("X-Real-IP", tt.ip)

			ctx.Request = &http.Request{RemoteAddr: tt.ip, Header: headers}
			m(ctx)

			if ctx.Writer.Status() != tt.want {
				t.Errorf("status = %v, want %v", ctx.Writer.Status(), tt.want)
			}
		})
	}
}

func TestCheckIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		subnets string
		want    bool
	}{
		{
			name:    "IP_not_exists",
			ip:      "127.0.32.1",
			subnets: "127.0.0.1/24",
			want:    false,
		},
		{
			name:    "IP_exists",
			ip:      "127.0.0.5",
			subnets: "127.0.0.1/24",
			want:    true,
		},
		{
			name:    "IP_empty",
			ip:      "",
			subnets: "127.0.0.1/24",
			want:    false,
		},
		{
			name:    "IP_many_subnet",
			ip:      "127.0.30.33",
			subnets: "127.0.0.1/24,127.0.30.1/24",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIP(tt.ip, tt.subnets); got != tt.want {
				t.Errorf("checkIP = %v, want %v", got, tt.want)
			}
		})
	}
}
