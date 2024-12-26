package service

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
)

type PingService struct {
	DatabaseDSN string
	log         *zap.Logger
}

func NewPingService(DatabaseDSN string, logger *zap.Logger) *PingService {
	return &PingService{DatabaseDSN: DatabaseDSN, log: logger}
}

func (s *PingService) Ping() error {
	db, err := sqlx.Connect("postgres", s.DatabaseDSN)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}

	return nil
}
