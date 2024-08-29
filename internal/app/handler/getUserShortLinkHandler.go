package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sanbright/go_shortener/internal/app/dto/user"
	"sanbright/go_shortener/internal/app/service"
)

type GetUserShortLinkHandler struct {
	service *service.ReadShortLinkService
	logger  *zap.Logger
}

func NewGetUserShortLinkHandler(service *service.ReadShortLinkService, logger *zap.Logger) *GetUserShortLinkHandler {
	return &GetUserShortLinkHandler{service: service, logger: logger}
}

func (handler *GetUserShortLinkHandler) Handle(ctx *gin.Context) {
	userId, ok := ctx.Get("UserId")
	if !ok {
		ctx.String(http.StatusUnauthorized, "")
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "application/json")
	defer ctx.Request.Body.Close()
	var res user.Response

	if str, ok := userId.(string); ok {

		shortLinksEntity, err := handler.service.GetByUserId(str)

		if err != nil || shortLinksEntity == nil {
			handler.logger.Info("Error get user short link",
				zap.Error(err),
			)
			ctx.String(http.StatusOK, "[]")
			ctx.Abort()
			return
		}

		for _, v := range *shortLinksEntity {
			res = append(res, &user.ItemResponse{
				OriginalUrl: v.URL,
				ShortURL:    v.ShortLink,
			})
		}
	}

	resp, _ := json.Marshal(res)

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(resp)

}
