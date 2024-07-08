package handler

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"strings"
	"testing"
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
	handler := NewPostShortLinkHandler(service.NewShortLinkService(shortLinkRepository, shortLinkGenerator))

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
			name:    "Success Append ShortLink",
			method:  http.MethodPost,
			request: "/",
			body:    "https://google.com/test",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/QYsTVwgznh",
			},
		},
		{
			name:    "Method invalid",
			method:  http.MethodGet,
			request: "/testesttest",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "Method not allowed!\n",
				location:   "",
			},
		},
		{
			name:    "Undefined Url",
			method:  http.MethodPost,
			request: "/testesttest",
			body:    "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "Not found url\n",
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			handler.Handle(response, request)

			result := response.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			body, err := io.ReadAll(result.Body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
