// Package handler обработчики запросов
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeleteUserShortLinkHandler обработчик удаления запроса на удаление коротких ссылок
type DeleteUserShortLinkHandler struct {
	service *service.WriteShortLinkService
	logger  *zap.Logger
	baseURL string
}

// NewDeleteUserShortLinkHandler конструктор обработчика удаления коротких ссылок
func NewDeleteUserShortLinkHandler(service *service.WriteShortLinkService, baseURL string, logger *zap.Logger) *DeleteUserShortLinkHandler {
	return &DeleteUserShortLinkHandler{service: service, logger: logger, baseURL: baseURL}
}

// Handle обработчик удаления запроса на удаление коротких ссылок
func (handler *DeleteUserShortLinkHandler) Handle(ctx *gin.Context) {
	userID, ok := ctx.Get("UserID")
	if !ok {
		ctx.String(http.StatusUnauthorized, "")
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "application/json")
	defer ctx.Request.Body.Close()

	var shortLinks []string

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &shortLinks)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	go handler.service.MarkAsRemove(shortLinks, userID.(string))

	ctx.Status(http.StatusAccepted)
}
