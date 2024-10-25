package main

import (
	"log"
	"os"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/middleware"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"sanbright/go_shortener/internal/config"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gin-contrib/pprof"
)

// Настройки по умолчанию
const (
	ShortLinkLen int    = 10
	CryptoKey    string = "$$ecuRityKe453H@"
)

func setupRouter(log *zap.Logger) *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(
		gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)),
		middleware.Logger(log),
		gin.Recovery(),
	)

	pprof.Register(r)

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
	logger := setupLogger()

	configuration, err := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"), os.Getenv("FILE_STORAGE_PATH"), os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("Fatal configuration error: %s", err.Error())
	}

	shortLinkGenerator := generator.NewShortLinkGenerator(ShortLinkLen)
	shortLinkRepository, _ := repository.NewRepositoryResolver(configuration, logger).Execute()

	readShortLinkService := service.NewReadShortLinkService(shortLinkRepository)
	writeShortLinkService := service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, logger)
	getHandler := handler.NewGetShortLinkHandler(readShortLinkService)
	postHandler := handler.NewPostShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL)
	postAPIHandler := handler.NewPostAPIShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL, logger)
	batchAPIHandler := handler.NewPostBatchShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL, logger)
	getPing := handler.NewGetPingHandler(configuration)
	getUserShortLinkHandler := handler.NewGetUserShortLinkHandler(readShortLinkService, configuration.BaseURL.URL, logger)
	deleteUserShortLinkHandler := handler.NewDeleteUserShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL, logger)

	cry := generator.NewCryptGenerator(CryptoKey)
	authMiddleware := middleware.Auth(cry, logger)
	authGenMiddleware := middleware.AuthGen(cry, configuration.DomainAndPort.Domain, logger)

	r := setupRouter(logger)
	r.GET(`/:id`, getHandler.Handle)
	r.POST(`/`, authGenMiddleware, postHandler.Handle)
	r.POST(`/api/shorten`, authGenMiddleware, postAPIHandler.Handle)
	r.POST(`/api/shorten/batch`, authGenMiddleware, batchAPIHandler.Handle)
	r.GET(`/api/user/urls`, authMiddleware, getUserShortLinkHandler.Handle)
	r.DELETE(`/api/user/urls`, authMiddleware, deleteUserShortLinkHandler.Handle)
	r.GET(`/ping`, getPing.Handle)
	err = r.Run(configuration.DomainAndPort.String())

	if err != nil {
		log.Fatalf("Fatal error: %s", err.Error())
	}
}
