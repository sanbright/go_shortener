package handler

import (
	"log"
	"net/http"
	"sanbright/go_shortener/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// GetPingHandler обработчик для проверки связи с базой данных
type GetPingHandler struct {
	Config *config.Config
}

// NewGetPingHandler конструктор обработчика проверки запроса на соединение с бд
func NewGetPingHandler(config *config.Config) *GetPingHandler {
	return &GetPingHandler{
		Config: config,
	}
}

// Handle выполнение запроса на проверку связи с бд
func (handler *GetPingHandler) Handle(ctx *gin.Context) {

	db, err := sqlx.Connect("postgres", handler.Config.DatabaseDSN)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("%s", err.Error())
		ctx.Abort()
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}
