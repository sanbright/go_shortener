package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sanbright/go_shortener/internal/app/service"
	"strings"
)

type PostShortLinkHandler struct {
	service *service.ShortLinkService
}

func NewPostShortLinkHandler(service *service.ShortLinkService) *PostShortLinkHandler {
	return &PostShortLinkHandler{service: service}
}

func (handler *PostShortLinkHandler) Handle(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodPost {
		ctx.String(http.StatusMethodNotAllowed, "Method not allowed!")
		ctx.Abort()
		return
	}

	uri := strings.TrimLeft(ctx.Request.RequestURI, "/")
	if len(uri) > 0 {
		ctx.String(http.StatusBadRequest, "Not found url")
		return
	}

	url, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	shortLinkEntity, err := handler.service.Add(string(url))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "text/plain")
	ctx.String(http.StatusCreated, "http://%s/%s", ctx.Request.Host, shortLinkEntity.ShortLink)

	return
}
