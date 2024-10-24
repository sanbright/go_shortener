package service

import (
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/repository"

	"github.com/google/uuid"
)

type ReadShortLinkService struct {
	repository repository.IShortLinkRepository
}

func NewReadShortLinkService(repository repository.IShortLinkRepository) *ReadShortLinkService {
	return &ReadShortLinkService{repository: repository}
}

func (service *ReadShortLinkService) GetByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	return service.repository.FindByShortLink(shortLink)
}

func (service *ReadShortLinkService) GetByUserID(userID string) (*[]entity.ShortLinkEntity, error) {
	return service.repository.FindByUserID(uuid.MustParse(userID))
}
