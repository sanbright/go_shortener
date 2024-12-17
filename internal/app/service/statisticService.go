package service

import (
	"go.uber.org/zap"
)

type StatisticService struct {
	service *ReadShortLinkService
	log     *zap.Logger
}

func NewStatisticService(service *ReadShortLinkService, logger *zap.Logger) *StatisticService {
	return &StatisticService{service: service, log: logger}
}

func (s *StatisticService) GetUrlsAndUsersStat() (int32, int32) {
	res, _ := s.service.GetStat()

	return int32(res.Urls), int32(res.Users)
}
