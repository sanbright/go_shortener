package handler

import (
	"errors"
	"io"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"strings"

	"net/http"

	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
)

// PostShortLinkHandler обработчик создания которкой ссылки
type PostShortLinkHandler struct {
	service *service.WriteShortLinkService
	baseURL string
}

// NewPostShortLinkHandler конструктор обработчик создания которкой ссылки
func NewPostShortLinkHandler(service *service.WriteShortLinkService, baseURL string) *PostShortLinkHandler {
	return &PostShortLinkHandler{service: service, baseURL: baseURL}
}

// Handle обработка запроса на создание короткой ссылки
func (handler *PostShortLinkHandler) Handle(ctx *gin.Context) {
	uri := strings.TrimLeft(ctx.Request.RequestURI, "/")
	if len(uri) > 0 {
		ctx.String(http.StatusNotFound, "Not found url")
		return
	}

	url, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	defer ctx.Request.Body.Close()

	userIDParam, ok := ctx.Get("UserID")
	if !ok {
		ctx.String(http.StatusUnauthorized, "")
		ctx.Abort()
		return
	}

	userID, _ := userIDParam.(string)

	shortLinkEntity, err := handler.service.Add(string(url), userID)
	statusCode := http.StatusCreated

	if err != nil {
		var notUniq *repErr.NotUniqShortLinkError

		if errors.As(err, &notUniq) {
			statusCode = http.StatusConflict
		} else {
			ctx.String(http.StatusBadRequest, "%s", err.Error())
			ctx.Abort()
			return
		}
	}

	ctx.Header("Content-type", "text/plain")
	ctx.String(statusCode, "%s/%s", handler.baseURL, shortLinkEntity.ShortLink)
}
