// Package service пакет для управления данными по коротким ссылкам
package service

import (
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/repository"

	"github.com/google/uuid"
)

// ReadShortLinkService - сервис для чтения данных по коротким ссылокам
type ReadShortLinkService struct {
	repository repository.IShortLinkRepository
}

// NewReadShortLinkService - конеструктор сервиса для чтения данных по коротким ссылокам
func NewReadShortLinkService(repository repository.IShortLinkRepository) *ReadShortLinkService {
	return &ReadShortLinkService{repository: repository}
}

// GetByShortLink - получение данных по краткой сслыке по короткой ссылке
// shortLink - короткая ссылка
func (service *ReadShortLinkService) GetByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	return service.repository.FindByShortLink(shortLink)
}

// GetByUserID - получение списка кратких ссылк пользователя
// userID - уникальный идентификатор пользователя
func (service *ReadShortLinkService) GetByUserID(userID string) (*[]entity.ShortLinkEntity, error) {
	return service.repository.FindByUserID(uuid.MustParse(userID))
}
