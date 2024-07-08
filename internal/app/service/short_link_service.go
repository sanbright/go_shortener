package service

import (
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/repository"
)

type ShortLinkService struct {
	repository repository.ShortLinkRepositoryInterface
	generator  generator.ShortLinkGeneratorInterface
}

func NewShortLinkService(repository repository.ShortLinkRepositoryInterface, generator generator.ShortLinkGeneratorInterface) *ShortLinkService {
	return &ShortLinkService{repository: repository, generator: generator}
}

func (service *ShortLinkService) GetByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	return service.repository.FindByShortLink(shortLink)
}

func (service *ShortLinkService) Add(url string) (*entity.ShortLinkEntity, error) {
	shortLink := service.generator.UniqGenerate()

	shortLinkEntity, err := service.repository.FindByShortLink(shortLink)

	if err == nil && shortLinkEntity != nil {
		shortLink = service.generator.UniqGenerate()
	}

	shortLinkEntity, err = service.repository.Add(shortLink, url)

	if err != nil {
		return nil, err
	}

	return shortLinkEntity, nil
}
