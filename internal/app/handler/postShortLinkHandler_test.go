package handler

import (
	"io"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/middleware"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
)

type MockShortLinkGenerator struct {
}

func NewMockShortLinkGenerator() *MockShortLinkGenerator {
	return &MockShortLinkGenerator{}
}

func (generator *MockShortLinkGenerator) UniqGenerate() string {
	return "QYsTVwgznh"
}

func TestPostShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := NewMockShortLinkGenerator()
	logger := setupLogger()
	handler := NewPostShortLinkHandler(service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, logger), "http://example.com")
	cry := generator.NewCryptGenerator("$$ecuRityKe453H@")
	authMiddleware := middleware.AuthGen(cry, "localhost", logger)

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
			name:    "UsageMethodNotAllowed",
			method:  http.MethodGet,
			request: "/",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "405 method not allowed",
				location:   "",
			},
		},
		{
			name:    "SuccessAppendShortLink",
			method:  http.MethodPost,
			request: "/",
			body:    "https://google.com/test",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/QYsTVwgznh",
			},
		},
		{
			name:    "UndefinedURL",
			method:  http.MethodPost,
			request: "/testesttest",
			body:    "",
			want: want{
				statusCode: http.StatusNotFound,
				body:       "404 page not found",
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			context, _ := gin.CreateTestContext(response)
			context.AddParam("short", strings.TrimLeft(tt.request, "/"))
			context.Request = request

			r := setupRouter()
			r.POST(`/`, authMiddleware, handler.Handle)
			r.ServeHTTP(response, request)
			result := response.Result()

			if code := tt.want.statusCode; code != result.StatusCode {
				t.Errorf("%v: StatusCode = '%v', want = '%v'", tt.name, result.StatusCode, code)
			}

			body, err := io.ReadAll(result.Body)
			defer result.Body.Close()

			if err != nil {
				t.Errorf("%v: Error = '%v'", tt.name, err.Error())
			}

			if tb := tt.want.body; tb != string(body) {
				t.Errorf("%v: Content = '%v', want = '%v'", tt.name, tb, string(body))
			}
		})
	}
}
