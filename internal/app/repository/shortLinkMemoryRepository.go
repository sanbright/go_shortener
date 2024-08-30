package repository

import (
	"fmt"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	repErr "sanbright/go_shortener/internal/app/repository/error"
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

func (repo *ShortLinkMemoryRepository) FindByURL(URL string) (*entity.ShortLinkEntity, error) {
	for _, v := range repo.Items {
		if v.URL == URL {
			return v, nil
		}
	}

	return nil, fmt.Errorf("not found by URL link: %s", URL)
}

func (repo *ShortLinkMemoryRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {
	found, _ := repo.FindByURL(url)

	if found != nil {
		return nil, repErr.NewNotUniqShortLinkError(found.URL, nil)
	}

	repo.Items[shortLink] = &entity.ShortLinkEntity{ShortLink: shortLink, URL: url}

	return repo.Items[shortLink], nil
}

func (repo *ShortLinkMemoryRepository) AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error) {
	for _, v := range shortLinks {
		_, err := repo.Add(v.ShortURL, v.OriginalURL)

		if err != nil {
			return nil, err
		}
	}

	return &shortLinks, nil
}
