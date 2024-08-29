package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sanbright/go_shortener/internal/app/dto/api"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
)

type PostAPIShortLinkHandler struct {
	service *service.WriteShortLinkService
	baseURL string
}

func NewPostAPIShortLinkHandler(service *service.WriteShortLinkService, baseURL string) *PostAPIShortLinkHandler {
	return &PostAPIShortLinkHandler{service: service, baseURL: baseURL}
}

func (handler *PostAPIShortLinkHandler) Handle(ctx *gin.Context) {

	var req *api.Request
	var buf bytes.Buffer

	_, err := buf.ReadFrom(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	if len(req.URL) == 0 {
		var out []*api.CurrentError
		out = append(out, &api.CurrentError{
			Path:    "url",
			Message: "Значение не может быть пустым",
		})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Success: false, Errors: out})
		return
	}

	userIdParam, ok := ctx.Get("UserId")
	if !ok {
		ctx.String(http.StatusUnauthorized, "")
		ctx.Abort()
		return
	}

	userId, ok := userIdParam.(string)

	shortLinkEntity, err := handler.service.Add(req.URL, userId)

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

	ctx.Header("Content-type", "application/json")

	res := api.Response{Result: handler.baseURL + "/" + shortLinkEntity.ShortLink}

	resp, _ := json.Marshal(res)

	ctx.Writer.WriteHeader(statusCode)
	ctx.Writer.Write(resp)
}
