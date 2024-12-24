package proto

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	repErr "sanbright/go_shortener/internal/app/repository/error"
	"sanbright/go_shortener/internal/app/service"

	"sanbright/go_shortener/internal/app/generator"
)

type IPingService interface {
	Ping() error
}

type IStatService interface {
	GetUrlsAndUsersStat() (int32, int32)
}

type IReadShortService interface {
	GetByUserID(context.Context, *PingRequest) (int32, int32)
}

type Server struct {
	UnimplementedServiceServer
	pingService IPingService
	statService IStatService
	generator   *generator.CryptGenerator
	rslService  *service.ReadShortLinkService
	wslService  *service.WriteShortLinkService
	baseUrl     string
}

func NewGPRCServer(pingService IPingService, statService IStatService, rslService *service.ReadShortLinkService, wslService *service.WriteShortLinkService, generator *generator.CryptGenerator, baseUrl string) *Server {
	return &Server{
		pingService: pingService,
		statService: statService,
		rslService:  rslService,
		wslService:  wslService,
		generator:   generator,
		baseUrl:     baseUrl,
	}
}

func (s *Server) GetStat(ctx context.Context, r *StatisticRequest) (*StatisticResponse, error) {
	var response StatisticResponse

	urlCount, userCount := s.statService.GetUrlsAndUsersStat()

	response.Urls = urlCount
	response.Users = userCount
	response.Code = http.StatusOK

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

func (s *Server) GetUsersURLs(ctx context.Context, r *GetByUserRequest) (*GetByUserResponse, error) {
	UUID, err := s.generator.DecodeValue(r.GetAuth())

	if err != nil {
		return nil, err
	}

	var response *GetByUserResponse
	var urls []*UserURL

	links, err := s.rslService.GetByUserID(UUID)

	if err != nil {
		return &GetByUserResponse{
			Code: http.StatusNoContent,
		}, err
	}

	for _, v := range *links {
		urls = append(urls, &UserURL{
			ShortUrl:    s.baseUrl + "/" + v.ShortLink,
			OriginalUrl: v.URL,
		})
	}

	statusCode := http.StatusOK

	if len(*links) == 0 {
		statusCode = http.StatusNoContent
	}

	response = &GetByUserResponse{
		Code: int32(statusCode),
		Urls: urls,
	}

	return response, nil
}

func (s *Server) PostShortLink(ctx context.Context, r *PostShortLinkRequest) (*PostShortLinkResponse, error) {
	userID, err := s.generator.DecodeValue(r.GetAuth())

	if err != nil {
		return nil, err
	}

	shortLinkEntity, err := s.wslService.Add(r.GetUrl(), userID)
	statusCode := http.StatusCreated

	if err != nil {
		var notUniq *repErr.NotUniqShortLinkError

		if errors.As(err, &notUniq) {
			statusCode = http.StatusConflict
		} else {
			statusCode = http.StatusBadRequest
			return &PostShortLinkResponse{
				Code: int32(statusCode),
			}, err
		}
	}

	return &PostShortLinkResponse{
		ShortUrl: fmt.Sprintf("%s/%s", s.baseUrl, shortLinkEntity.ShortLink),
		Code:     int32(statusCode),
	}, nil
}

func (s *Server) DeleteURLs(ctx context.Context, r *DeleteRequest) (*DeleteResponse, error) {
	userID, err := s.generator.DecodeValue(r.GetAuth())

	if err != nil {
		return &DeleteResponse{Code: int32(http.StatusBadRequest)}, nil
	}

	go s.wslService.MarkAsRemove(r.GetUrls(), userID)

	return &DeleteResponse{Code: int32(http.StatusAccepted)}, nil
}
