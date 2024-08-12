package repository

import (
	"fmt"
	"sanbright/go_shortener/internal/app/entity"
)

type ShortLinkMemoryRepository struct {
	Items map[string]*entity.ShortLinkEntity
}

func NewShortLinkRepository() *ShortLinkMemoryRepository {
	return &ShortLinkMemoryRepository{
		Items: make(map[string]*entity.ShortLinkEntity),
	}
}

func (repo *ShortLinkMemoryRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	if shortLinkEntity, exists := repo.Items[shortLink]; exists {
		return shortLinkEntity, nil
	}

	return nil, fmt.Errorf("not found by short link: %s", shortLink)
}

func (repo *ShortLinkMemoryRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {

	repo.Items[shortLink] = &entity.ShortLinkEntity{ShortLink: shortLink, URL: url}

	return repo.Items[shortLink], nil
}
