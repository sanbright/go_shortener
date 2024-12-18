package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
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
