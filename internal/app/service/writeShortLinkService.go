package service

import (
	"errors"
	"sanbright/go_shortener/internal/app/dto/batch"
	"sanbright/go_shortener/internal/app/entity"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/repository"
	repErr "sanbright/go_shortener/internal/app/repository/error"
)

type WriteShortLinkService struct {
	repository repository.ShortLinkRepositoryInterface
	generator  generator.ShortLinkGeneratorInterface
}

func NewWriteShortLinkService(repository repository.ShortLinkRepositoryInterface, generator generator.ShortLinkGeneratorInterface) *WriteShortLinkService {
	return &WriteShortLinkService{repository: repository, generator: generator}
}

func (service *WriteShortLinkService) Add(url string, userId string) (*entity.ShortLinkEntity, error) {
	shortLink := service.generator.UniqGenerate()

	shortLinkEntity, err := service.repository.Add(shortLink, url, userId)

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

func (service *WriteShortLinkService) AddBatch(links *batch.Request, userId string) (*batch.AddBatchDtoList, error) {
	var batchList batch.AddBatchDtoList

	for _, element := range *links {
		batchList = append(batchList, &batch.AddBatchDto{
			CorrelationID: element.CorrelationID,
			OriginalURL:   element.OriginalURL,
			ShortURL:      service.generator.UniqGenerate(),
			UserId:        userId,
		})
	}

	result, err := service.repository.AddBatch(batchList)

	if err != nil {
		return &batchList, err
	}

	return result, nil
}
