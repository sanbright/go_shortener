package repository

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	repErr "sanbright/go_shortener/internal/app/repository/error"
)

type ShortLinkDBRepository struct {
	db *sqlx.DB
}

func NewShortLinkDBRepository(db *sqlx.DB) *ShortLinkDBRepository {
	return &ShortLinkDBRepository{db: db}
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

	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(pgerrcode.UniqueViolation) {
			return nil, repErr.NewNotUniqShortLinkError(newShortLinkEntity.URL, err)
		}

		return nil, err
	}

	return newShortLinkEntity, nil
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

func (repo *ShortLinkDBRepository) FindByURL(URL string) (*entity.ShortLinkEntity, error) {
	var shortLinkEntity entity.ShortLinkEntity

	err := repo.db.Get(&shortLinkEntity,
		`SELECT 
 					uuid,
					short_link,
					url
				FROM short_link sl
					WHERE sl.url = $1
				LIMIT 1`,
		URL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &shortLinkEntity, err
}

func (repo *ShortLinkDBRepository) AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error) {
	tx, err := repo.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	for _, v := range shortLinks {

		im := entity.NewShortLinkEntity(v.ShortURL, v.OriginalURL)

		_, err := tx.Exec(
			"INSERT INTO short_link (uuid, short_link, url) VALUES ($1, $2, $3)",
			im.UUID,
			im.ShortLink,
			im.URL)

		if err != nil {
			if pgerrcode.IsIntegrityConstraintViolation(pgerrcode.UniqueViolation) {
				return nil, repErr.NewNotUniqShortLinkError(v.OriginalURL, err)
			}

			return nil, err
		}
	}

	tx.Commit()

	return &shortLinks, nil
}
