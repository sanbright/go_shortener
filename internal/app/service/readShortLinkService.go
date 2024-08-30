package service

import (
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/repository"
)

type ReadShortLinkService struct {
	repository repository.ShortLinkRepositoryInterface
}

func NewReadShortLinkService(repository repository.ShortLinkRepositoryInterface) *ReadShortLinkService {
	return &ReadShortLinkService{repository: repository}
}

func (service *ReadShortLinkService) GetByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	return service.repository.FindByShortLink(shortLink)
}
