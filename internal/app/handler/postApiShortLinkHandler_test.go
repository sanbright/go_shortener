package handler

import (
	"io"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
)

func TestPostApiShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := NewMockShortLinkGenerator()
	handler := NewPostApiShortLinkHandler(service.NewShortLinkService(shortLinkRepository, shortLinkGenerator), "http://example.com")

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
			request: "/api/shorten",
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
			request: "/api/shorten",
			body:    "{\"url\":\"https://google.com/test\"}",
			want: want{
				statusCode: http.StatusCreated,
				body:       "{\"result\":\"http://example.com/QYsTVwgznh\"}",
			},
		},
		{
			name:    "InvalidUrl",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    "{\"urasdl\":\"https://google.com/test\"}",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"success\":false,\"errors\":[{\"path\":\"url\",\"message\":\"Значение не может быть пустым\"}]}",
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

			context.Writer.Header().Set("Accept-Encoding", "gzip")
			context.Writer.Header().Set("Content-Encoding", "gzip")
			context.Writer.Header().Set("Content-Type", "application/json")

			r := setupRouter()
			r.POST(`/api/shorten`, handler.Handle)
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
