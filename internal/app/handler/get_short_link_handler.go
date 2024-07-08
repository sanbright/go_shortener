package handler

import (
	"net/http"
	"sanbright/go_shortener/internal/app/service"
	"strings"
)

type GetShortLinkHandler struct {
	service *service.ShortLinkService
}

func NewGetShortLinkHandler(service *service.ShortLinkService) *GetShortLinkHandler {
	return &GetShortLinkHandler{service: service}
}

func (handler *GetShortLinkHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	shortLinkEntity, err := handler.service.GetByShortLink(strings.TrimLeft(request.RequestURI, "/"))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(writer, request, shortLinkEntity.Url, http.StatusTemporaryRedirect)

	return
}
