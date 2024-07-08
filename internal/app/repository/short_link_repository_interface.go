package repository

import (
	"sanbright/go_shortener/internal/app/entity"
)

type ShortLinkRepositoryInterface interface {
	FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error)
	Add(shortLink string, url string) (*entity.ShortLinkEntity, error)
}
