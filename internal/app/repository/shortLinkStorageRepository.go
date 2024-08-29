package repository

import (
	"bufio"
	"github.com/google/uuid"
	"io"
	"os"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"strings"
)

import (
	"encoding/json"
)

type ShortLinkStorageRepository struct {
	file *os.File
}

func NewShortLinkStorageRepository(path string) (*ShortLinkStorageRepository, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortLinkStorageRepository{file: file}, nil
}

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

func (repo *ShortLinkStorageRepository) FindByUserId(uuid uuid.UUID) (*[]entity.ShortLinkEntity, error) {
	var entityList *[]entity.ShortLinkEntity
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

			*entityList = append(*entityList, shortLinkEntity)
		}
	}

	return entityList, nil
}

func (repo *ShortLinkStorageRepository) Add(shortLink string, url string, userId string) (*entity.ShortLinkEntity, error) {
	found, _ := repo.FindByURL(url)

	if found != nil {
		return nil, repErr.NewNotUniqShortLinkError(found.URL, nil)
	}

	var newShortLinkEntity = entity.NewShortLinkEntity(shortLink, url, userId)

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

func (repo *ShortLinkStorageRepository) AddBatch(shortLinks batch.AddBatchDtoList) (*batch.AddBatchDtoList, error) {
	for _, v := range shortLinks {
		_, err := repo.Add(v.ShortURL, v.OriginalURL, v.UserId)

		if err != nil {
			return nil, err
		}
	}

	return &shortLinks, nil
}
