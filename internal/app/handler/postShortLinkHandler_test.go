package handler

import (
	"io"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/middleware"
	"sanbright/go_shortener/internal/app/proto"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"

	"github.com/gin-gonic/gin"
)

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
			name:    "ConflictAppendShortLink",
			method:  http.MethodPost,
			request: "/",
			body:    "https://google.com/test",
			want: want{
				statusCode: http.StatusConflict,
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

func TestPostShortLinkHandler_GRPC(t *testing.T) {
	type want struct {
		statusCode int
		body       string
		location   string
	}

	tests := []struct {
		name string
		auth string
		body string
		want want
	}{
		{
			name: "SuccessAppendShortLink",
			auth: "1FLRobWnu0pYInXBHnJmU8T3GOvB86FawJeOUdZDVYB7ido58lc8mLgBXaUzKAoydcxieg==",
			body: "https://google.com/test",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/QYsTVwgznh",
			},
		},
		{
			name: "ConflictAppendShortLink",
			auth: "1FLRobWnu0pYInXBHnJmU8T3GOvB86FawJeOUdZDVYB7ido58lc8mLgBXaUzKAoydcxieg==",
			body: "https://google.com/test",
			want: want{
				statusCode: http.StatusConflict,
				body:       "http://example.com/QYsTVwgznh",
			},
		},
	}

	client, ctx, conn := setupGRPCClient()

	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.PostShortLink(ctx, &proto.PostShortLinkRequest{Url: tt.body, Auth: tt.auth})

			if err != nil {
				t.Fatalf("PostShortLink failed: %v", err)
			}

			if code := tt.want.statusCode; code != int(response.Code) {
				t.Errorf("%v: StatusCode = '%v', want = '%v'", tt.name, code, int(response.Code))
			}

			if url := tt.want.body; url != response.ShortUrl {
				t.Errorf("%v: ShortUrl = '%v', want = '%v'", tt.name, url, response.ShortUrl)
			}
		})
	}
}
