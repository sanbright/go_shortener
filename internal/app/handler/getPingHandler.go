package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sanbright/go_shortener/internal/app/service"
	// asd
	_ "github.com/lib/pq"
)

// GetPingHandler обработчик для проверки связи с базой данных
type GetPingHandler struct {
	pingService *service.PingService
}

// NewGetPingHandler конструктор обработчика проверки запроса на соединение с бд
func NewGetPingHandler(pingService *service.PingService) *GetPingHandler {
	return &GetPingHandler{
		pingService: pingService,
	}
}

// Handle выполнение запроса на проверку связи с бд
func (handler *GetPingHandler) Handle(ctx *gin.Context) {

	err := handler.pingService.Ping()

	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}
