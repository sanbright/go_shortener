package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/config"
)

const schema = `
CREATE TABLE IF NOT EXISTS short_link (
	"uuid" UUID NOT NULL,
	"short_link" VARCHAR(10) NOT NULL,
	"url" TEXT NOT NULL,
	PRIMARY KEY ("uuid")
);

CREATE UNIQUE INDEX IF NOT EXISTS short_link__uniq ON short_link (short_link);
CREATE UNIQUE INDEX IF NOT EXISTS url__uniq ON short_link (url);
`

type ShortLinkRepositoryInterface interface {
	Add(shortLink string, url string) (*entity.ShortLinkEntity, error)
	AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error)
	FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error)
	FindByURL(URL string) (*entity.ShortLinkEntity, error)
}

type Resolver struct {
	Config *config.Config
	Log    *zap.Logger
}

func NewRepositoryResolver(config *config.Config, log *zap.Logger) *Resolver {
	return &Resolver{Config: config, Log: log}
}

func (r *Resolver) Execute() (ShortLinkRepositoryInterface, error) {
	if len(r.Config.DatabaseDSN) > 0 {
		db, err := sqlx.Connect("postgres", r.Config.DatabaseDSN)
		if err != nil {
			r.Log.Error("Fatal init DB repository:", zap.String("ERROR", err.Error()))
		}
		if db != nil {
			err = db.Ping()

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

func (r *Resolver) InitDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", r.Config.DatabaseDSN)
	if err != nil {
		r.Log.Error("Fatal init DB repository:", zap.String("ERROR", err.Error()))
		return nil, err
	}

	return db, nil
}