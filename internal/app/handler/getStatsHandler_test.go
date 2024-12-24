package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/proto"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"strings"
	"testing"
)

func TestGetStatsHandler_Handle(t *testing.T) {
	shortLinkRepository := repository.NewShortLinkRepository()

	_, err := shortLinkRepository.Add("sa42d45ds2", "https:\\\\testing.com\\ksjadkjas", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		t.Errorf("ShortLinkFixture: Error = '%v'", err.Error())
	}

	_, err = shortLinkRepository.Add("qwetyr123iu", "https:\\\\google.com", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	handler := NewGetStatsHandler(service.NewReadShortLinkService(shortLinkRepository))

	type want struct {
		statusCode int
		body       string
		location   string
	}

	tests := []struct {
		name        string
		method      string
		contentType string
		request     string
		body        string
		want        want
	}{
		{
			name:    "SuccessGetStatsHandler_1",
			method:  http.MethodGet,
			request: "/api/internal/stats",
			want: want{
				statusCode: http.StatusOK,
				body:       "{\"urls\":2,\"users\":1}",
			},
		},
		{
			name:    "UsageMethodNotAllowed",
			method:  http.MethodPost,
			request: "/api/internal/stats",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "405 method not allowed",
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))

			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.Request = request

			r := setupRouter()
			r.GET(`/api/internal/stats`, handler.Handle)
			r.ServeHTTP(response, request)

			if code := tt.want.statusCode; code != response.Code {
				t.Errorf("%v: StatusCode = '%v', want = '%v'", tt.name, response.Code, code)
			}

			body, err := io.ReadAll(response.Body)

			if err != nil {
				t.Errorf("%v: Error = '%v'", tt.name, err.Error())
			}

			if tbody := tt.want.body; tbody != string(body) {
				t.Errorf("%v: Content = '%v', want = '%v'", tt.name, tbody, string(body))
			}
		})
	}
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	log := setupLogger()

	shortLinkRepository := repository.NewShortLinkRepository()

	_, err := shortLinkRepository.Add("sa42d45ds2", "https:\\\\testing.com\\ksjadkjas", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		log.Error("ShortLinkFixture: Error = '%v'", zap.Error(err))
	}

	_, err = shortLinkRepository.Add("qwetyr123iu", "https:\\\\google.com", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		log.Error("hortLinkFixture: Error = '%v'", zap.Error(err))
	}

	shortLinkGenerator := NewMockShortLinkGenerator()
	writeShortLinkService := service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, log)
	readShortLinkService := service.NewReadShortLinkService(shortLinkRepository)

	pingService := service.NewPingService("configuration.DatabaseDSN", log)
	statService := service.NewStatisticService(readShortLinkService, log)
	cry := generator.NewCryptGenerator("$$ecuRityKe453H@")
	proto.RegisterServiceServer(s, proto.NewGPRCServer(pingService, statService, readShortLinkService, writeShortLinkService, cry, "http://example.com"))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Error("Server exited with error: %v", zap.Error(err))
		}
	}()
}

func TestGetStatsHandler_GRPC(t *testing.T) {

	type want struct {
		statusCode int32
		urls       int32
		users      int32
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "SuccessGetStatsGRPC_1",
			want: want{
				statusCode: http.StatusOK,
				users:      1,
				urls:       2,
			},
		},
	}

	client, ctx, conn := setupGRPCClient()

	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.GetStat(ctx, &proto.StatisticRequest{})

			if err != nil {
				t.Fatalf("SayHello failed: %v", err)
			}
			log.Printf("GetStat: %+v", response)

			if code := tt.want.statusCode; code != response.Code {
				t.Errorf("%v: StatusCode = '%v', want = '%v'", tt.name, response.Code, code)
			}

			if urls := tt.want.urls; urls != response.Urls {
				t.Errorf("%v: Content = '%v', want = '%v'", tt.name, urls, string(response.Urls))
			}

			if users := tt.want.users; users != response.Users {
				t.Errorf("%v: Content = '%v', want = '%v'", tt.name, users, string(response.Users))
			}
		})
	}

}
