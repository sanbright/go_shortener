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

type ReadShortLinkRepository struct {
	file *os.File
}

type WriteShortLinkRepository struct {
	file *os.File
}

func NewWriteShortLinkRepository(path string) (*WriteShortLinkRepository, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &WriteShortLinkRepository{file: file}, nil
}

func NewReadShortLinkRepository(path string) (*ReadShortLinkRepository, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ReadShortLinkRepository{file: file}, nil
}

func (repo *ReadShortLinkRepository) FindByShortLink(shortLink string) (*entity.ShortLinkEntity, error) {
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

func (repo *WriteShortLinkRepository) Add(shortLink string, url string) (*entity.ShortLinkEntity, error) {
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
