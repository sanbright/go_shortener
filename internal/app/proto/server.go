package proto

import (
	"context"
)

type IPingService interface {
	Ping() error
}

type IStatService interface {
	GetUrlsAndUsersStat() (int32, int32)
}

type Server struct {
	UnimplementedServiceServer
	pingService IPingService
	statService IStatService
}

func NewGPRCServer(pingService IPingService, statService IStatService) *Server {
	return &Server{
		pingService: pingService,
		statService: statService,
	}
}

func (s *Server) GetStat(ctx context.Context, r *StatisticRequest) (*StatisticResponse, error) {
	var response StatisticResponse

	urlCount, userCount := s.statService.GetUrlsAndUsersStat()

	response.Urls = urlCount
	response.Users = userCount

	return &response, nil
}

func (s *Server) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	var response PingResponse

	err := s.pingService.Ping()
	if err != nil {
		response.Error = err.Error()
	}

	return &response, nil
}
