package handler

import (
	"encoding/json"
	"net/http"
	"sanbright/go_shortener/internal/app/dto/user"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUserShortLinkHandler обработчик получение списка которких ссылкок созданных пользователем
type GetUserShortLinkHandler struct {
	service *service.ReadShortLinkService
	logger  *zap.Logger
	baseURL string
}

// NewGetUserShortLinkHandler конструктор обработчика получения пользовательских ссылок
func NewGetUserShortLinkHandler(service *service.ReadShortLinkService, baseURL string, logger *zap.Logger) *GetUserShortLinkHandler {
	return &GetUserShortLinkHandler{service: service, logger: logger, baseURL: baseURL}
}

// Handle обработка запроса на получение списка ссылок для пользователя
func (handler *GetUserShortLinkHandler) Handle(ctx *gin.Context) {
	userID, ok := ctx.Get("UserID")
	if !ok {
		ctx.String(http.StatusUnauthorized, "")
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "application/json")
	defer ctx.Request.Body.Close()
	var res user.Response

	if str, ok := userID.(string); ok {

		shortLinksEntity, err := handler.service.GetByUserID(str)

		if err != nil || shortLinksEntity == nil || len(*shortLinksEntity) == 0 {
			handler.logger.Info("Error get user short link",
				zap.Error(err),
			)

			ctx.String(http.StatusNoContent, "[]")
			ctx.Abort()
			return
		}

		for _, v := range *shortLinksEntity {
			res = append(res, &user.ItemResponse{
				OriginalURL: v.URL,
				ShortURL:    handler.baseURL + "/" + v.ShortLink,
			})
		}
	}

	resp, _ := json.Marshal(res)

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(resp)

}
