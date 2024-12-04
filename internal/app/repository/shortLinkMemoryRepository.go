package repository

import (
	"fmt"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	repErr "sanbright/go_shortener/internal/app/repository/error"

	"github.com/google/uuid"
)

// ShortLinkMemoryRepository хранилище коротких ссылкок в памяти
type ShortLinkMemoryRepository struct {
	Items map[string]*entity.ShortLinkEntity
}

// NewShortLinkRepository конструктор хранилища коротких ссылкок в памяти
func NewShortLinkRepository() *ShortLinkMemoryRepository {
	return &ShortLinkMemoryRepository{
		Items: make(map[string]*entity.ShortLinkEntity),
	}
}

// FindByShortLink происк в хранилище информации по короткой ссылке
//
//	shortLink - краткая ссылка
func (repo *ShortLinkMemoryRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	if shortLinkEntity, exists := repo.Items[shortLink]; exists {
		return shortLinkEntity, nil
	}

	return nil, fmt.Errorf("not found by short link: %s", shortLink)
}

// FindByURL происк в хранилище информации по ссылке
//
//	URL - оригнальная ссылка
func (repo *ShortLinkMemoryRepository) FindByURL(URL string) (*entity.ShortLinkEntity, error) {
	for _, v := range repo.Items {
		if v.URL == URL {
			return v, nil
		}
	}

	return nil, fmt.Errorf("not found by URL link: %s", URL)
}

// FindByUserID получение списка коротких ссылок из хранилища для конкретного пользователя
//
//	uuid - уникальный идентификатор пользователя
func (repo *ShortLinkMemoryRepository) FindByUserID(uuid uuid.UUID) (*[]entity.ShortLinkEntity, error) {
	var entityList []entity.ShortLinkEntity

	for _, v := range repo.Items {
		if v.UserID == uuid {
			entityList = append(entityList, *v)
		}
	}

	return &entityList, nil
}

// Add добавление информации по короткой ссылке в хранилище
//
//	shortLink - краткая ссылка
//	url - оригинальный URL
//	userID - UUID пользователя
func (repo *ShortLinkMemoryRepository) Add(shortLink string, url string, userID string) (*entity.ShortLinkEntity, error) {
	found, _ := repo.FindByURL(url)

	if found != nil {
		return nil, repErr.NewNotUniqShortLinkError(found.URL, nil)
	}

	repo.Items[shortLink] = entity.NewShortLinkEntity(shortLink, url, userID)

	return repo.Items[shortLink], nil
}

// AddBatch добавление пачки коротких ссылок.
//
//	shortLinks - список добавляемых коротких ссылок.
func (repo *ShortLinkMemoryRepository) AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error) {
	for _, v := range shortLinks {
		_, err := repo.Add(v.ShortURL, v.OriginalURL, v.UserID)

		if err != nil {
			return nil, err
		}
	}

	return &shortLinks, nil
}

// Delete удаление списка коротких ссылок.
//
//		shortLinkList - список удаляемых коротких ссылок.
//	 userID - идентификатор пользователя
func (repo *ShortLinkMemoryRepository) Delete(shortLinkList []string, userID string) error {
	for _, shortLink := range shortLinkList {
		repo.Items[shortLink].IsDeleted = true
	}

	return nil
}

// GetStat получение статистики по коротким ссылкам
func (repo *ShortLinkMemoryRepository) GetStat() (int, int, error) {
	return 0, 0, nil
}
