package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/middleware"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"strings"
	"testing"
)

func TestPostBathShortLinkHandler_Handle(t *testing.T) {
	shortLinkRepository := repository.NewShortLinkRepository()
	shortLinkGenerator := NewMockShortLinkGenerator()
	cry := generator.NewCryptGenerator("$$ecuRityKe453H@")

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	handler := NewPostBatchShortLinkHandler(service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, logger), "http://example.com", logger)
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
			request: "/api/shorten/batch",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "405 method not allowed",
				location:   "",
			},
		},
		{
			name:    "SuccessAppendOneBatchShortLink",
			method:  http.MethodPost,
			request: "/api/shorten/batch",
			body:    "[{\"correlation_id\":\"asdasdas\",\"original_url\":\"https://google.com\"}]",
			want: want{
				statusCode: http.StatusCreated,
				body:       "[{\"correlation_id\":\"asdasdas\",\"short_url\":\"http://example.com/QYsTVwgznh\"}]",
			},
		},
		{
			name:    "SuccessAppendTwoBatchShortLink",
			method:  http.MethodPost,
			request: "/api/shorten/batch",
			body:    "[{\"correlation_id\":\"dd112935-bb0f-4645-bb19-49a1418ba692\",\"original_url\":\"http://tdbju0uipsn1ng.com/xirnkpgj9cvjyn/x2jdb5h9ltw\"},{\"correlation_id\":\"2cfbbb87-643c-4dfa-a4c4-a40b9213188f\",\"original_url\":\"http://mpxjbc26zly.biz/vdhyungaq6m\"}]",
			want: want{
				statusCode: http.StatusCreated,
				body:       "[{\"correlation_id\":\"dd112935-bb0f-4645-bb19-49a1418ba692\",\"short_url\":\"http://example.com/QYsTVwgznh\"},{\"correlation_id\":\"2cfbbb87-643c-4dfa-a4c4-a40b9213188f\",\"short_url\":\"http://example.com/QYsTVwgznh\"}]",
			},
		},
		{
			name:    "ConflictAppendTwoBatchShortLink",
			method:  http.MethodPost,
			request: "/api/shorten/batch",
			body:    "[{\"correlation_id\":\"dd112935-bb0f-4645-bb19-49a1418ba692\",\"original_url\":\"http://tdbju0uipsn1ng.com/xirnkpgj9cvjyn/x2jdb5h9ltw\"},{\"correlation_id\":\"2cfbbb87-643c-4dfa-a4c4-a40b9213188f\",\"original_url\":\"http://tdbju0uipsn1ng.com/xirnkpgj9cvjyn/x2jdb5h9ltw\"}]",
			want: want{
				statusCode: http.StatusConflict,
				body:       "[{\"correlation_id\":\"dd112935-bb0f-4645-bb19-49a1418ba692\",\"short_url\":\"http://example.com/QYsTVwgznh\"},{\"correlation_id\":\"2cfbbb87-643c-4dfa-a4c4-a40b9213188f\",\"short_url\":\"http://example.com/QYsTVwgznh\"}]",
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
			r.POST(`/api/shorten/batch`, authMiddleware, handler.Handle)
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
