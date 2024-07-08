package handler

import (
	"io"
	"net/http"
	"sanbright/go_shortener/internal/app/service"
	"strings"
)

type PostShortLinkHandler struct {
	service *service.ShortLinkService
}

func NewPostShortLinkHandler(service *service.ShortLinkService) *PostShortLinkHandler {
	return &PostShortLinkHandler{service: service}
}

func (handler *PostShortLinkHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	uri := strings.TrimLeft(request.RequestURI, "/")
	if len(uri) > 0 {
		http.Error(writer, "Not found url", http.StatusBadRequest)
		return
	}

	url, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortLinkEntity, err := handler.service.Add(string(url))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	if _, err = writer.Write([]byte("http://" + request.Host + "/" + shortLinkEntity.ShortLink)); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-type", "text/plain")

	return
}
