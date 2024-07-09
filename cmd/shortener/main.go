package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanbright/go_shortener/internal/app/generator"
	handler2 "sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"sanbright/go_shortener/internal/config"
)

func main() {
	configuration := config.NewConfig()
	flag.Var(&configuration.DomainAndPort, "a", "listen host and port")
	flag.Var(&configuration.BaseUrl, "b", "domain in short link")
	flag.Parse()

	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := generator.NewShortLinkGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortLinkGenerator)
	getHandler := handler2.NewGetShortLinkHandler(shortLinkService)
	postHandler := handler2.NewPostShortLinkHandler(shortLinkService, configuration.BaseUrl.URL)

	ginRouter := gin.Default()
	ginRouter.POST(`/`, postHandler.Handle)
	ginRouter.GET(`/:id`, getHandler.Handle)

	err := http.ListenAndServe(configuration.DomainAndPort.String(), ginRouter)

	if err != nil {
		panic(err)
	}
}
