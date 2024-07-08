package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sanbright/go_shortener/internal/app/service"
)

type GetShortLinkHandler struct {
	service *service.ShortLinkService
}

func NewGetShortLinkHandler(service *service.ShortLinkService) *GetShortLinkHandler {
	return &GetShortLinkHandler{service: service}
}

func (handler *GetShortLinkHandler) Handle(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		ctx.String(http.StatusMethodNotAllowed, "Method not allowed!")
		ctx.Abort()
		return
	}

	shortLinkEntity, err := handler.service.GetByShortLink(ctx.Param("id"))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, shortLinkEntity.Url)

	return
}
