package repository

import (
	"fmt"
	"sanbright/go_shortener/internal/app/entity"
)

type ShortLinkRepository struct {
	Items map[string]string
}

func NewShortLinkRepository() *ShortLinkRepository {
	return &ShortLinkRepository{
		Items: make(map[string]string),
	}
}

func (repo *ShortLinkRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	if url, exists := repo.Items[shortLink]; exists {
		return &(entity.ShortLinkEntity{ShortLink: shortLink, Url: url}), nil
	}

	return nil, fmt.Errorf("not found by short link: %s", shortLink)
}

func (repo *ShortLinkRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {
	repo.Items[shortLink] = url

	return &entity.ShortLinkEntity{ShortLink: shortLink, Url: url}, nil
}
