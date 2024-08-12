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
	configuration, err := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"), os.Getenv("FILE_STORAGE_PATH"), os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("Fatal configuration error: %s", err.Error())
	}

	shortLinkGenerator := generator.NewShortLinkGenerator(ShortLinkLen)
	writeShortLinkRepository, err := repository.NewWriteShortLinkRepository(configuration.StoragePath)
	if err != nil {
		log.Fatalf("Fatal init write repository: %s", err.Error())
	}

	readShortLinkRepository, _ := repository.NewReadShortLinkRepository(configuration.StoragePath)
	if err != nil {
		log.Fatalf("Fatal init read repository: %s", err.Error())
	}

	readShortLinkService := service.NewReadShortLinkService(readShortLinkRepository)
	writeShortLinkService := service.NewWriteShortLinkService(writeShortLinkRepository, shortLinkGenerator)
	getHandler := handler.NewGetShortLinkHandler(readShortLinkService)
	postHandler := handler.NewPostShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL)
	postAPIHandler := handler.NewPostAPIShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL)
	getPing := handler.NewGetPingHandler(configuration)

	r := setupRouter()
	r.GET(`/:id`, getHandler.Handle)
	r.POST(`/`, postHandler.Handle)
	r.POST(`/api/shorten`, postAPIHandler.Handle)
	r.GET(`/ping`, getPing.Handle)
	err = r.Run(configuration.DomainAndPort.String())

	if err != nil {
		log.Fatalf("Fatal error: %s", err.Error())
	}
}
