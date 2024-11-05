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

func TestDeleteUserShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()

	_, err := shortLinkRepository.Add("iqwnmqw9001", "https:\\\\testing.com\\ksjadkjas", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		t.Errorf("ShortLinkFixture: Error = '%v'", err.Error())
	}

	_, err = shortLinkRepository.Add("hj3393893fn", "https:\\\\google.com", "653e7307-6960-4b60-ab1b-44cd2f662634")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	logger := setupLogger()
	shortLinkGenerator := NewMockShortLinkGenerator()
	handler := NewDeleteUserShortLinkHandler(service.NewWriteShortLinkService(shortLinkRepository, shortLinkGenerator, logger), "http://example.com", logger)
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
			name:    "SuccessRemoveUserShortLink_1",
			method:  http.MethodDelete,
			auth:    "I8LumVeMYJlq8pNoeeY0s1EzbMS90OFaFnH0uXYKv3I7FEbDBSDPMvRjLDgVZx3Q8wGSGA==",
			request: "/api/user/urls",
			body:    "[\"iqwnmqw9001\"]",
			want: want{
				statusCode: http.StatusAccepted,
				body:       "",
			},
		},
		{
			name:    "SuccessRemoveUserShortLink_2",
			method:  http.MethodDelete,
			auth:    "1FLRobWnu0pYInXBHnJmU8T3GOvB86FawJeOUdZDVYBg+invalid==",
			request: "/api/user/urls",
			body:    "[\"iqwnmqw9001\",\"hj3393893fn\"]",
			want: want{
				statusCode: http.StatusAccepted,
				body:       "",
			},
		},
		{
			name:    "FailRemoveUserShortLink",
			method:  http.MethodDelete,
			auth:    "1FLRobWnu0pYInXBHnJmU8T3GOvB86FawJeOUdZDVYBg+invalid==",
			request: "/api/user/urls",
			body:    "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"Offset\":0}",
			},
		},
		{
			name:    "FailJSONRemoveUserShortLink",
			method:  http.MethodDelete,
			auth:    "1FLRobWnu0pYInXBHnJmU8T3GOvB86FawJeOUdZDVYBg+invalid==2",
			request: "/api/user/urls",
			body:    "{{",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"Offset\":2}",
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
		{
			name:    "UsageMethodNotAllowed",
			method:  http.MethodGet,
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
			r.DELETE(`/api/user/urls`, authMiddleware, handler.Handle)
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
