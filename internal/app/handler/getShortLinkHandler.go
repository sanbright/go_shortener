package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sanbright/go_shortener/internal/app/service"
)

type GetShortLinkHandler struct {
	service *service.ReadShortLinkService
}

func NewGetShortLinkHandler(service *service.ReadShortLinkService) *GetShortLinkHandler {
	return &GetShortLinkHandler{service: service}
}

func (handler *GetShortLinkHandler) Handle(ctx *gin.Context) {
	shortLinkEntity, err := handler.service.GetByShortLink(ctx.Param("id"))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	if shortLinkEntity == nil {
		ctx.String(http.StatusNotFound, "Not found link")
		ctx.Abort()
		return
	}

	defer ctx.Request.Body.Close()

	ctx.Redirect(http.StatusTemporaryRedirect, shortLinkEntity.URL)
}
