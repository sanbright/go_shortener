package repository

import (
	"bufio"
	"io"
	"os"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"strings"

	"github.com/google/uuid"
)

import (
	"encoding/json"
)

// ShortLinkStorageRepository файловое хранилище данных по коротким ссылкам
type ShortLinkStorageRepository struct {
	file *os.File
}

// NewShortLinkStorageRepository конструктор файлового хранилища данных по коротким ссылкам
func NewShortLinkStorageRepository(path string) (*ShortLinkStorageRepository, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortLinkStorageRepository{file: file}, nil
}

// FindByShortLink происк в хранилище информации по короткой ссылке
//
//	shortLink - краткая ссылка
func (repo *ShortLinkStorageRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
	var shortLinkEntity entity.ShortLinkEntity

	_, err := repo.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(repo.file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), shortLink) {
			err := json.Unmarshal(scanner.Bytes(), &shortLinkEntity)
			if err != nil {
				return nil, err
			}

			return &shortLinkEntity, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

// FindByURL происк в хранилище информации по ссылке
//
//	URL - оригнальная ссылка
func (repo *ShortLinkStorageRepository) FindByURL(URL string) (*entity.ShortLinkEntity, error) {
	var shortLinkEntity entity.ShortLinkEntity

	_, err := repo.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(repo.file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), URL) {
			err := json.Unmarshal(scanner.Bytes(), &shortLinkEntity)
			if err != nil {
				return nil, err
			}

			return &shortLinkEntity, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

// FindByUserID получение списка коротких ссылок из хранилища для конкретного пользователя
//
//	uuid - уникальный идентификатор пользователя
func (repo *ShortLinkStorageRepository) FindByUserID(uuid uuid.UUID) (*[]entity.ShortLinkEntity, error) {
	var entityList []entity.ShortLinkEntity
	var shortLinkEntity entity.ShortLinkEntity

	_, err := repo.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(repo.file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), uuid.String()) {
			err := json.Unmarshal(scanner.Bytes(), &shortLinkEntity)
			if err != nil {
				return nil, err
			}

			entityList = append(entityList, shortLinkEntity)
		}
	}

	return &entityList, nil
}

// Add добавление информации по короткой ссылке в хранилище
//
//	shortLink - краткая ссылка
//	url - оригинальный URL
//	userID - UUID пользователя
func (repo *ShortLinkStorageRepository) Add(shortLink string, url string, userID string) (*entity.ShortLinkEntity, error) {
	found, _ := repo.FindByURL(url)

	if found != nil {
		return nil, repErr.NewNotUniqShortLinkError(found.URL, nil)
	}

	var newShortLinkEntity = entity.NewShortLinkEntity(shortLink, url, userID)

	s, err := json.Marshal(newShortLinkEntity)

	if err != nil {
		return nil, err
	}

	_, err = repo.file.WriteString(string(s) + "\n")

	if err != nil {
		return nil, err
	}

	return newShortLinkEntity, nil
}

// AddBatch добавление пачки коротких ссылок.
//
//	shortLinks - список добавляемых коротких ссылок.
func (repo *ShortLinkStorageRepository) AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error) {
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
func (repo *ShortLinkStorageRepository) Delete(shortLinkList []string, userID string) error {
	return nil
}
