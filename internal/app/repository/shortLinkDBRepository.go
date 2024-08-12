package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sanbright/go_shortener/internal/app/entity"
)

type ShortLinkDBRepository struct {
	db *sqlx.DB
}

func NewShortLinkDBRepository(db *sqlx.DB) *ShortLinkDBRepository {
	return &ShortLinkDBRepository{db: db}
}

func (repo *ShortLinkDBRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	var shortLinkEntity entity.ShortLinkEntity

	err := repo.db.Get(&shortLinkEntity,
		`SELECT 
 					uuid,
					short_link,
					url
				FROM short_link sl
					WHERE sl.short_link = $1
				LIMIT 1`,
		shortLink)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &shortLinkEntity, err
}

func (repo *ShortLinkDBRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {
	var newShortLinkEntity = entity.NewShortLinkEntity(shortLink, url)

	_, err := repo.db.Exec(
		"INSERT INTO short_link (uuid, short_link, url) VALUES ($1, $2, $3)",
		newShortLinkEntity.UUID,
		newShortLinkEntity.ShortLink,
		newShortLinkEntity.URL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return newShortLinkEntity, nil
}
