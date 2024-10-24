package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sanbright/go_shortener/internal/app/dto/batch"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostBatchShortLinkHandler struct {
	service *service.WriteShortLinkService
	logger  *zap.Logger
	baseURL string
}

func NewPostBatchShortLinkHandler(service *service.WriteShortLinkService, baseURL string, logger *zap.Logger) *PostBatchShortLinkHandler {
	return &PostBatchShortLinkHandler{service: service, baseURL: baseURL, logger: logger}
}

func (handler *PostBatchShortLinkHandler) Handle(ctx *gin.Context) {
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
	ctx.Header("Content-type", "application/json")

	userIDParam, ok := ctx.Get("UserID")
	if !ok {
		ctx.String(http.StatusUnauthorized, "")
		ctx.Abort()
		return
	}

	userID, _ := userIDParam.(string)

	list, err := handler.service.AddBatch(req, userID)

	statusCode := http.StatusCreated

	if err != nil {
		handler.logger.Error("Add Batch Error", zap.Error(err))
		var notUniq *repErr.NotUniqShortLinkError

		if errors.As(err, &notUniq) {
			statusCode = http.StatusConflict
			handler.logger.Error("Add Batch Conflict")
		} else {
			ctx.String(http.StatusBadRequest, "%s", err.Error())
			ctx.Abort()
			return
		}
	}

	for _, element := range *list {
		out = append(out, &batch.ItemResponse{
			CorrelationID: element.CorrelationID,
			ShortURL:      handler.baseURL + "/" + element.ShortURL,
		})
	}

	resp, _ := json.Marshal(out)

	ctx.Writer.WriteHeader(statusCode)
	ctx.Writer.Write(resp)
}
