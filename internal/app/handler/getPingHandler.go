package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sanbright/go_shortener/internal/config"
)

type GetPingHandler struct {
	Config *config.Config
}

func NewGetPingHandler(config *config.Config) *GetPingHandler {
	return &GetPingHandler{
		Config: config,
	}
}
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
	return
}
