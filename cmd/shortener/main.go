package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/handler"
	"sanbright/go_shortener/internal/app/middleware"
	"sanbright/go_shortener/internal/app/proto"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"sanbright/go_shortener/internal/config"
	"syscall"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gin-contrib/pprof"
	"google.golang.org/grpc"
)

// BuildVersion = определяет версию приложения
// BuildDate = определяет дату сборки
// BuildCommit = определяет коммит сборки
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
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
	configuration, err := config.NewConfig(
		os.Getenv("SERVER_ADDRESS"),
		os.Getenv("BASE_URL"),
		os.Getenv("FILE_STORAGE_PATH"),
		os.Getenv("DATABASE_DSN"),
		os.Getenv("ENABLE_HTTPS") == "true",
		os.Getenv("CONFIG"),
		os.Getenv("TRUSTED_SUBNET"),
		os.Getenv("GRPC_HOST"),
	)

	if err != nil {
		log.Fatalf("Fatal configuration error: %s", err.Error())
	}

	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	configuration.GRPCHost = ":8088"

	fmt.Printf("GRPCHost: %s\n", configuration.GRPCHost)
	logger := setupLogger()

	if configuration.GRPCHost != "" {
		initGRPCServer(configuration, logger)
	} else {
		r := initServer(configuration, logger)

		srv := &http.Server{
			Addr:    configuration.DomainAndPort.String(),
			Handler: r,
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			var srvErr error
			if configuration.HTTPS {
				fmt.Printf("SSL mode\n")
				srvErr = srv.ListenAndServeTLS("./key.crt", "./key.pem")
			} else {
				fmt.Printf("not SSL mode\n")
				srvErr = srv.ListenAndServe()
			}

			if srvErr != nil {
				log.Printf("Fatal error: %s", srvErr.Error())
			}

			os.Exit(0)
		}()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-stop

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to gracefully shutdown server: %s", err.Error())
		}

		cancel()
	}
}

func initService(configuration *config.Config, log *zap.Logger) (*service.PingService, *service.StatisticService, *service.ReadShortLinkService, *service.WriteShortLinkService, *generator.CryptGenerator) {
	shortLinkGenerator := generator.NewShortLinkGenerator(ShortLinkLen)

	shortLinkRepository, _ := repository.NewRepositoryResolver(configuration, log).Execute()
	readShortLinkService := service.NewReadShortLinkService(shortLinkRepository)
	writeShortLinkService := service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, log)

	return service.NewPingService(configuration.DatabaseDSN, log),
		service.NewStatisticService(readShortLinkService, log),
		readShortLinkService,
		writeShortLinkService,
		generator.NewCryptGenerator(CryptoKey)
}

func initGRPCServer(configuration *config.Config, log *zap.Logger) {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", configuration.GRPCHost)
	if err != nil {
		log.Error("Add Batch Error", zap.Error(err))
	}

	s := grpc.NewServer()
	pingService, statService, readShortLinkService, writeShortLinkService, cry := initService(configuration, log)

	proto.RegisterServiceServer(s, proto.NewGPRCServer(pingService, statService, readShortLinkService, writeShortLinkService, cry, configuration.BaseURL.String()))

	fmt.Println("Сервер gRPC начал работу " + configuration.GRPCHost)
	// получаем запрос gRPC
	if err = s.Serve(listen); err != nil {
		log.Error("Failure listen", zap.Error(err))
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Failure recover", zap.Any("error", err))
		}
	}()
}

func initServer(configuration *config.Config, log *zap.Logger) *gin.Engine {
	pingService, _, _, _, cry := initService(configuration, log)

	shortLinkGenerator := generator.NewShortLinkGenerator(ShortLinkLen)
	shortLinkRepository, _ := repository.NewRepositoryResolver(configuration, log).Execute()

	readShortLinkService := service.NewReadShortLinkService(shortLinkRepository)
	writeShortLinkService := service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, log)
	getHandler := handler.NewGetShortLinkHandler(readShortLinkService)
	postHandler := handler.NewPostShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL)
	postAPIHandler := handler.NewPostAPIShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL, log)
	batchAPIHandler := handler.NewPostBatchShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL, log)

	getPing := handler.NewGetPingHandler(pingService)
	getUserShortLinkHandler := handler.NewGetUserShortLinkHandler(readShortLinkService, configuration.BaseURL.URL, log)
	deleteUserShortLinkHandler := handler.NewDeleteUserShortLinkHandler(writeShortLinkService, configuration.BaseURL.URL, log)
	getStats := handler.NewGetStatsHandler(readShortLinkService)

	authMiddleware := middleware.Auth(cry, log)
	authGenMiddleware := middleware.AuthGen(cry, configuration.DomainAndPort.Domain, log)
	trustedSubnetMiddleware := middleware.IpFilter(log, configuration.TrustedSubnet)

	r := setupRouter(log)
	r.GET(`/:id`, getHandler.Handle)
	r.POST(`/`, authGenMiddleware, postHandler.Handle)
	r.POST(`/api/shorten`, authGenMiddleware, postAPIHandler.Handle)
	r.POST(`/api/shorten/batch`, authGenMiddleware, batchAPIHandler.Handle)
	r.GET(`/api/user/urls`, authMiddleware, getUserShortLinkHandler.Handle)
	r.DELETE(`/api/user/urls`, authMiddleware, deleteUserShortLinkHandler.Handle)
	r.GET(`/ping`, getPing.Handle)
	r.GET(`/api/internal/stats`, trustedSubnetMiddleware, getStats.Handle)

	return r
}
