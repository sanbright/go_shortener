package service

import (
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/repository"
)

type WriteShortLinkService struct {
	repository repository.WriteShortLinkRepositoryInterface
	generator  generator.ShortLinkGeneratorInterface
}

func NewWriteShortLinkService(repository repository.WriteShortLinkRepositoryInterface, generator generator.ShortLinkGeneratorInterface) *WriteShortLinkService {
	return &WriteShortLinkService{repository: repository, generator: generator}
}

func (service *WriteShortLinkService) Add(url string) (*entity.ShortLinkEntity, error) {
	shortLink := service.generator.UniqGenerate()

	shortLinkEntity, err := service.repository.Add(shortLink, url)

	if err != nil {
		return nil, err
	}

	return shortLinkEntity, nil
}
