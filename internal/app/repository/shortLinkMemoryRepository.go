package repository

import (
	"fmt"
	"sanbright/go_shortener/internal/app/entity"
)

type WriteShortLinkRepositoryInterface interface {
	Add(shortLink string, url string) (*entity.ShortLinkEntity, error)
}

type ReadShortLinkRepositoryInterface interface {
	FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error)
}

type ShortLinkRepository struct {
	Items map[string]*entity.ShortLinkEntity
}

func NewShortLinkRepository() *ShortLinkRepository {
	return &ShortLinkRepository{
		Items: make(map[string]*entity.ShortLinkEntity),
	}
}

func (repo *ShortLinkRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	if shortLinkEntity, exists := repo.Items[shortLink]; exists {
		return shortLinkEntity, nil
	}

	return nil, fmt.Errorf("not found by short link: %s", shortLink)
}

func (repo *ShortLinkRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {

	repo.Items[shortLink] = &entity.ShortLinkEntity{ShortLink: shortLink, URL: url}

	return repo.Items[shortLink], nil
}
