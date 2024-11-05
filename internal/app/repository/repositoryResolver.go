// Package repository пакет для работы с хранилищем данных
package repository

import (
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/config"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const schema = `
CREATE TABLE IF NOT EXISTS short_link (
	"uuid" UUID NOT NULL,
	"short_link" VARCHAR(10) NOT NULL,
	"url" TEXT NOT NULL,
	"user_id" UUID,
	"is_deleted" BOOLEAN NOT NULL DEFAULT false,
	PRIMARY KEY ("uuid")
);

CREATE UNIQUE INDEX IF NOT EXISTS short_link__uniq ON short_link (short_link);
CREATE UNIQUE INDEX IF NOT EXISTS url__uniq ON short_link (url);
`

// IShortLinkRepository интерфейс для работы с хранилищем данных по коротким ссылкам
type IShortLinkRepository interface {
	Add(shortLink string, url string, userID string) (*entity.ShortLinkEntity, error)
	AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error)
	FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error)
	FindByURL(URL string) (*entity.ShortLinkEntity, error)
	FindByUserID(uuid uuid.UUID) (*[]entity.ShortLinkEntity, error)
	Delete(shortLinkList []string, userID string) error
}

// Resolver резолвер, определяет какое хранилище необходимо использовать и возвращает его для использования
type Resolver struct {
	// Config объект конфигурации проекта
	Config *config.Config
	// Log логгер
	Log *zap.Logger
	// DB подключение к СУБД
	DB *sqlx.DB
}

// NewRepositoryResolver Конструктор резолвера хранилища
// config - объект конфигурации проекта
// log логгер
func NewRepositoryResolver(config *config.Config, log *zap.Logger) *Resolver {
	return &Resolver{Config: config, Log: log}
}

// Execute - резолв хранилища
func (r *Resolver) Execute() (IShortLinkRepository, error) {
	if len(r.Config.DatabaseDSN) > 0 {
		db, _ := r.InitDB()

		if db != nil {
			err := db.Ping()

			if err != nil {
				r.Log.Error("Error ping DB repository:", zap.String("ERROR", err.Error()))
			}

			db.MustExec(schema)
			r.Log.Debug("Used DB repository")
			return NewShortLinkDBRepository(db), nil
		}
	}

	if len(r.Config.StoragePath) > 0 {
		storageRepository, err := NewShortLinkStorageRepository(r.Config.StoragePath)
		if err != nil {
			r.Log.Error("Fatal init File repository:", zap.String("ERROR", err.Error()))
		}

		if storageRepository != nil {
			r.Log.Debug("Used File repository")
			return storageRepository, nil
		}
	}

	r.Log.Debug("Used Memory repository")
	return NewShortLinkRepository(), nil
}

// InitDB - установка соединения с СУБД
func (r *Resolver) InitDB() (*sqlx.DB, error) {
	if r.DB != nil {
		return r.DB, nil
	}

	db, err := sqlx.Connect("postgres", r.Config.DatabaseDSN)
	if err != nil {
		r.Log.Error("Fatal init DB repository:", zap.String("ERROR", err.Error()))
		return nil, err
	}

	return db, nil
}
