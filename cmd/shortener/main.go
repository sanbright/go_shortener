package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"sanbright/go_shortener/internal/config"
)

func main() {
	configuration := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"))
	flag.Var(&configuration.DomainAndPort, "a", "listen host and port")
	flag.Var(&configuration.BaseUrl, "b", "domain in short link")
	flag.Parse()

	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := generator.NewShortLinkGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortLinkGenerator)
	getHandler := handler.NewGetShortLinkHandler(shortLinkService)
	postHandler := handler.NewPostShortLinkHandler(shortLinkService, configuration.BaseUrl.URL)

	ginRouter := gin.Default()
	ginRouter.POST(`/`, postHandler.Handle)
	ginRouter.GET(`/:id`, getHandler.Handle)

	err := http.ListenAndServe(configuration.DomainAndPort.String(), ginRouter)

	if err != nil {
		panic(err)
	}
}
