package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/service"
)

type PostBathShortLinkHandler struct {
	service *service.WriteShortLinkService
	baseURL string
}

func NewPostBathShortLinkHandler(service *service.WriteShortLinkService, baseURL string) *PostBathShortLinkHandler {
	return &PostBathShortLinkHandler{service: service, baseURL: baseURL}
}

func (handler *PostBathShortLinkHandler) Handle(ctx *gin.Context) {
	var req *batch.Request
	var buf bytes.Buffer
	var out batch.Response

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

	for _, element := range *req {
		shortLinkEntity, _ := handler.service.Add(element.OriginalURL)
		out = append(out, &batch.ItemResponse{
			CorrelationID: element.CorrelationID,
			ShortURL:      handler.baseURL + "/" + shortLinkEntity.ShortLink,
		})
	}
	resp, _ := json.Marshal(out)

	ctx.Header("Content-type", "application/json")
	ctx.Writer.WriteHeader(http.StatusCreated)
	ctx.Writer.Write(resp)
}
