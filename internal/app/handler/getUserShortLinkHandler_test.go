package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/app/generator"
	"sanbright/go_shortener/internal/app/middleware"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUserShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()

	_, err := shortLinkRepository.Add("sa42d45ds2", "https:\\\\testing.com\\ksjadkjas", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		t.Errorf("ShortLinkFixture: Error = '%v'", err.Error())
	}

	_, err = shortLinkRepository.Add("qwetyr123iu", "https:\\\\google.com", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	logger := setupLogger()

	handler := NewGetUserShortLinkHandler(service.NewReadShortLinkService(shortLinkRepository), "http://example.com", logger)
	cry := generator.NewCryptGenerator("$$ecuRityKe453H@")

	authMiddleware := middleware.Auth(cry, logger)
	type want struct {
		statusCode int
		body       string
		location   string
	}

	tests := []struct {
		name        string
		method      string
		auth        string
		contentType string
		request     string
		body        string
		want        want
	}{
		{
			name:    "SuccessGettingUserShortLink_1",
			method:  http.MethodGet,
			auth:    "I8LumVeMYJlq8pNoeeY0s1EzbMS90OFaFnH0uXYKv3I7FEbDBSDPMvRjLDgVZx3Q8wGSGA==",
			request: "/api/user/urls",
			want: want{
				statusCode: http.StatusOK,
				body:       "[{\"original_url\":\"https:\\\\\\\\testing.com\\\\ksjadkjas\",\"short_url\":\"http://example.com/sa42d45ds2\"},{\"original_url\":\"https:\\\\\\\\google.com\",\"short_url\":\"http://example.com/qwetyr123iu\"}]",
			},
		},
		{
			name:    "SuccessGettingUserShortLink_2",
			method:  http.MethodGet,
			auth:    "1FLRobWnu0pYInXBHnJmU8T3GOvB86FawJeOUdZDVYB7ido58lc8mLgBXaUzKAoydcxieg==",
			request: "/api/user/urls",
			want: want{
				statusCode: http.StatusNoContent,
				body:       "",
			},
		},
		{
			name:    "UsageMethodNotAllowed",
			method:  http.MethodPost,
			request: "/api/user/urls",
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
			request.AddCookie(&http.Cookie{Name: "Auth", Value: tt.auth})

			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.AddParam("id", strings.TrimLeft(tt.request, "/"))
			context.AddParam("userId", "4c1b4334-8f1c-4874-8750-c5214e2f48b9")
			context.Request = request

			r := setupRouter()
			r.GET(`/api/user/urls`, authMiddleware, handler.Handle)
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
