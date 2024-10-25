// Package service пакет для управления данными по коротким ссылкам
package service

import (
	"errors"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/repository"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// WriteShortLinkService - сервис для записи коротких ссылок
type WriteShortLinkService struct {
	repository repository.IShortLinkRepository
	generator  generator.IShortLinkGenerator
	logger     *zap.Logger
}

// NewWriteShortLinkService - конструктор сервиса для записи коротких ссылок
func NewWriteShortLinkService(repository repository.IShortLinkRepository, generator generator.IShortLinkGenerator, logger *zap.Logger) *WriteShortLinkService {
	return &WriteShortLinkService{repository: repository, generator: generator, logger: logger}
}

// Add - добавление короткой ссылки
// url - оригинальный УРЛ
// userID - уникальный идентификатор пользователя
// возвращает сущность созданной короткой ссылки
func (service *WriteShortLinkService) Add(url string, userID string) (*entity.ShortLinkEntity, error) {
	shortLink := service.generator.UniqGenerate()

	shortLinkEntity, err := service.repository.Add(shortLink, url, userID)

	if err != nil {
		var notUniq *repErr.NotUniqShortLinkError

		if errors.As(err, &notUniq) {
			shortLinkEntity, _ = service.repository.FindByURL(url)

			return shortLinkEntity, err
		}

		return nil, err
	}

	return shortLinkEntity, nil
}

// AddBatch - добавление пачки коротких ссылок
// links - скисок коротких ссылок
// userID - уникальный идентификатор пользователя
// возвращает спискок коротких ссылок
func (service *WriteShortLinkService) AddBatch(links *batch.Request, userID string) (*batch.AddBatchDtoList, error) {
	var batchList batch.AddBatchDtoList

	for _, element := range *links {
		batchList = append(batchList, &batch.AddBatchDto{
			CorrelationID: element.CorrelationID,
			OriginalURL:   element.OriginalURL,
			ShortURL:      service.generator.UniqGenerate(),
			UserID:        userID,
		})
	}

	result, err := service.repository.AddBatch(batchList)

	if err != nil {
		return &batchList, err
	}

	return result, nil
}

// MarkAsRemove - пометить короткую ссылку как удаленную
// shortLinkList - скисок коротких ссылок
// userID - уникальный идентификатор пользователя
// возвращает список помеченных коротких ссылок
func (service *WriteShortLinkService) MarkAsRemove(shortLinkList []string, userID string) []string {
	var chunk []string
	var chunks [][]string
	i := 0
	for _, shortLink := range shortLinkList {
		chunk = append(chunk, shortLink)

		if i < 30 {
			i++
		} else {
			chunks = append(chunks, chunk)
			chunk = nil
			i = 0
		}

	}

	var deletedLinks []string
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
		chunk = nil
	}

	inCh := service.sendToPrepare(chunks)
	ch1 := service.prepareRemoveShortLink(inCh, userID)
	ch2 := service.prepareRemoveShortLink(inCh, userID)
	for n := range service.fanIn(ch1, ch2) {
		deletedLinks = append(deletedLinks, n...)
	}

	return deletedLinks
}

// sendToPrepare - отправка чанков на обработку
// chunks - скисок коротких ссылок
// возвращает список помеченных коротких ссылок
func (service *WriteShortLinkService) sendToPrepare(chunks [][]string) chan []string {
	outCh := make(chan []string)
	go func() {
		defer close(outCh)
		for _, chunk := range chunks {
			outCh <- chunk
		}
	}()

	return outCh
}

// prepareRemoveShortLink - отправка чанков на обработку
// inCh - масив чанков
// userID - уникальный идентификатор пользователя
// возвращает список помеченных коротких ссылок
func (service *WriteShortLinkService) prepareRemoveShortLink(inCh chan []string, userID string) chan []string {
	outCh := make(chan []string)

	go func() {
		defer close(outCh)
		for shortLinks := range inCh {
			err := service.repository.Delete(shortLinks, userID)
			if err != nil {
				service.logger.Error(
					"Ошибка удаления записей",
					zap.String("shortLinks", strings.Join(shortLinks, ",")),
					zap.String("userID", userID),
					zap.Error(err),
				)
				return
			}

			outCh <- shortLinks
		}
	}()

	return outCh
}

func (service *WriteShortLinkService) fanIn(chs ...chan []string) chan []string {
	var wg sync.WaitGroup
	outCh := make(chan []string)

	output := func(c chan []string) {
		for n := range c {
			outCh <- n
		}
		wg.Done()
	}

	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}
