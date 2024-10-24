package handler

import (
	"net/http"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
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

	if shortLinkEntity.IsDeleted {
		ctx.String(http.StatusGone, "Not found link")
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, shortLinkEntity.URL)
}
