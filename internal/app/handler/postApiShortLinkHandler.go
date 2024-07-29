package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sanbright/go_shortener/internal/app/dto"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
)

type PostApiShortLinkHandler struct {
	service *service.ShortLinkService
	baseURL string
}

func NewPostApiShortLinkHandler(service *service.ShortLinkService, baseURL string) *PostApiShortLinkHandler {
	return &PostApiShortLinkHandler{service: service, baseURL: baseURL}
}

func (handler *PostApiShortLinkHandler) Handle(ctx *gin.Context) {

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

	if len(req.Url) == 0 {
		var out []*dto.CurrentError
		out = append(out, &dto.CurrentError{
			Path:    "url",
			Message: "Значение не может быть пустым",
		})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Errors: out})
		return
	}

	shortLinkEntity, err := handler.service.Add(req.Url)

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
