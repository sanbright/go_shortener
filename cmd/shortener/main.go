package main

import (
	"github.com/gin-contrib/gzip"
	"log"
	"os"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/middleware"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"sanbright/go_shortener/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const ShortLinkLen int = 10

func setupRouter() *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(
		gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)),
		middleware.Logger(setupLogger()),
		gin.Recovery(),
	)

	return r
}

func setupLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}

func main() {
	configuration, err := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"))
	if err != nil {
		log.Fatalf("Fatal configuration error: %s", err.Error())
	}

	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := generator.NewShortLinkGenerator(ShortLinkLen)
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortLinkGenerator)
	getHandler := handler.NewGetShortLinkHandler(shortLinkService)
	postHandler := handler.NewPostShortLinkHandler(shortLinkService, configuration.BaseURL.URL)
	postApiHandler := handler.NewPostApiShortLinkHandler(shortLinkService, configuration.BaseURL.URL)

	r := setupRouter()
	r.GET(`/:id`, getHandler.Handle)
	r.POST(`/`, postHandler.Handle)
	r.POST(`/api/shorten`, postApiHandler.Handle)
	err = r.Run(configuration.DomainAndPort.String())

	if err != nil {
		log.Fatalf("Fatal error: %s", err.Error())
	}
}
