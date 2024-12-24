package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanbright/go_shortener/internal/app/service"
)

type GetStatsHandler struct {
	service *service.ReadShortLinkService
}

// NewGetStatsHandler конструктор обработчика запросов по статистике
func NewGetStatsHandler(service *service.ReadShortLinkService) *GetStatsHandler {
	return &GetStatsHandler{service: service}
}

// Handle выполнение запроса на получение коротких ссылок
func (handler *GetStatsHandler) Handle(ctx *gin.Context) {
	res, _ := handler.service.GetStat()

	resp, _ := json.Marshal(res)

	ctx.Header("Content-type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(resp)
}
