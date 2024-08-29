package service

import (
	"github.com/google/uuid"
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

func (service *ReadShortLinkService) GetByUserId(userId string) (*[]entity.ShortLinkEntity, error) {
	return service.repository.FindByUserId(uuid.MustParse(userId))
}
