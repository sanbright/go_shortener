package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sanbright/go_shortener/internal/app/generator"
	handler2 "sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
)

func main() {
	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := generator.NewShortLinkGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortLinkGenerator)
	getHandler := handler2.NewGetShortLinkHandler(shortLinkService)
	postHandler := handler2.NewPostShortLinkHandler(shortLinkService)

	ginRouter := gin.Default()
	ginRouter.POST(`/`, postHandler.Handle)
	ginRouter.GET(`/:id`, getHandler.Handle)

	err := http.ListenAndServe(`localhost:8080`, ginRouter)

	if err != nil {
		panic(err)
	}
}
