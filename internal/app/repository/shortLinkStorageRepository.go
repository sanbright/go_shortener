package repository

import (
	"bufio"
	"io"
	"os"
	"sanbright/go_shortener/internal/app/entity"
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

func (repo *ShortLinkStorageRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {
	var newShortLinkEntity = entity.NewShortLinkEntity(shortLink, url)

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
