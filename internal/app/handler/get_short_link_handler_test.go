package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/app/repository"
	"sanbright/go_shortener/internal/app/service"
	"strings"
	"testing"
)

func TestGetShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()

	_, _ = shortLinkRepository.Add("sa42d45ds2", "https:\\\\testing.com\\ksjadkjas")
	_, _ = shortLinkRepository.Add("qwetyr123iu", "https:\\\\google.com")
	shortLinkGenerator := NewMockShortLinkGenerator()
	handler := NewGetShortLinkHandler(service.NewShortLinkService(shortLinkRepository, shortLinkGenerator))

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
			name:    "Success Getting ShortLink",
			method:  http.MethodGet,
			request: "/sa42d45ds2",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"https:\\\\testing.com\\ksjadkjas\">Temporary Redirect</a>.\n\n",
				location:   "https:\\\\testing.com\\ksjadkjas",
			},
		},
		{
			name:    "Success Getting ShortLink",
			method:  http.MethodGet,
			request: "/qwetyr123iu",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"https:\\\\google.com\">Temporary Redirect</a>.\n\n",
				location:   "https:\\\\google.com",
			},
		},
		{
			name:    "Method invalid",
			method:  http.MethodPost,
			request: "/",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "Method not allowed!",
				location:   "",
			},
		},
		{
			name:    "Undefined Url",
			method:  http.MethodGet,
			request: "/testesttest",
			body:    "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "not found by short link: testesttest",
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))

			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.AddParam("id", strings.TrimLeft(tt.request, "/"))
			context.Request = request

			handler.Handle(context)

			assert.Equal(t, tt.want.statusCode, response.Code)

			body, err := io.ReadAll(response.Body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
			assert.Equal(t, response.Header().Get("Location"), tt.want.location)
		})
	}
}
