package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

type ShortLinkDBRepository struct {
	db *sqlx.DB
}

func NewShortLinkDBRepository(db *sqlx.DB) *ShortLinkDBRepository {
	return &ShortLinkDBRepository{db: db}
}

func (repo *ShortLinkDBRepository) Add(shortLink string, url string, userID string) (*entity.ShortLinkEntity, error) {
	var newShortLinkEntity = entity.NewShortLinkEntity(shortLink, url, userID)

	_, err := repo.db.Exec(
		"INSERT INTO short_link (uuid, short_link, url, user_id) VALUES ($1, $2, $3, $4)",
		newShortLinkEntity.UUID,
		newShortLinkEntity.ShortLink,
		newShortLinkEntity.URL,
		newShortLinkEntity.UserID,
	)

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
					url,
					user_id,
					is_deleted
				FROM short_link sl
					WHERE sl.short_link = $1
				LIMIT 1`,
		shortLink)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &shortLinkEntity, err
}

func (repo *ShortLinkDBRepository) FindByURL(URL string) (*entity.ShortLinkEntity, error) {
	var shortLinkEntity entity.ShortLinkEntity

	err := repo.db.Get(&shortLinkEntity,
		`SELECT 
 					uuid,
					short_link,
					url,
					user_id,
					is_deleted
				FROM short_link sl
					WHERE sl.url = $1
					AND sl.is_deleted = false
				LIMIT 1`,
		URL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &shortLinkEntity, err
}

func (repo *ShortLinkDBRepository) FindByUserID(uuid uuid.UUID) (*[]entity.ShortLinkEntity, error) {
	var shortLinkEntities []entity.ShortLinkEntity

	err := repo.db.Select(&shortLinkEntities,
		`SELECT 
 					uuid,
					short_link,
					url,
					user_id,
					is_deleted
				FROM short_link sl
				WHERE sl.user_id = $1
				AND sl.is_deleted = false`,
		uuid.String())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &shortLinkEntities, err
}

func (repo *ShortLinkDBRepository) AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error) {
	tx, err := repo.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	for _, v := range shortLinks {

		im := entity.NewShortLinkEntity(v.ShortURL, v.OriginalURL, v.UserID)

		_, err := tx.Exec(
			"INSERT INTO short_link (uuid, short_link, url, user_id) VALUES ($1, $2, $3, $4)",
			im.UUID,
			im.ShortLink,
			im.URL,
			im.UserID)

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

func (repo *ShortLinkDBRepository) Delete(shortLinkList []string, userID string) error {
	var inArray []string
	var params []interface{}

	params = append(params, userID)
	i := 2
	for _, shortLink := range shortLinkList {
		inArray = append(inArray, "$"+fmt.Sprintf("%d", i))
		params = append(params, shortLink)
		i++
	}
	_, err := repo.db.Exec(
		fmt.Sprintf(`update short_link SET is_deleted = true WHERE user_id = $1 AND short_link IN (%s)`, strings.Join(inArray, ",")), params...)

	if err != nil {
		return err
	}

	return nil
}
