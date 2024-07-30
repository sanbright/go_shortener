package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sanbright/go_shortener/internal/app/dto"
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

	var req *dto.Request
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
		var out []*dto.CurrentError
		out = append(out, &dto.CurrentError{
			Path:    "url",
			Message: "Значение не может быть пустым",
		})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Errors: out})
		return
	}

	shortLinkEntity, err := handler.service.Add(req.URL)

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "application/json")

	res := dto.Response{Result: handler.baseURL + "/" + shortLinkEntity.ShortLink}

	resp, _ := json.Marshal(res)

	ctx.Writer.WriteHeader(http.StatusCreated)
	ctx.Writer.Write(resp)
}
