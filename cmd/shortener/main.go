package main

import (
	"net/http"
	handler2 "sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
)

func main() {
	serveMux := http.NewServeMux()

	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkService := service.NewShortLinkService(shortLinkRepository)
	getHandler := handler2.NewGetShortLinkHandler(shortLinkService)
	postHandler := handler2.NewPostShortLinkHandler(shortLinkService)

	serveMux.HandleFunc(`/`, postHandler.Handle)
	serveMux.HandleFunc(`/{id}`, getHandler.Handle)

	err := http.ListenAndServe(`localhost:8083`, serveMux)

	if err != nil {
		panic(err)
	}
}
