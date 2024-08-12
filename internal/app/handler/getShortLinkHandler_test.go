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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.HandleMethodNotAllowed = true

	return r
}

func TestGetShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()

	_, err := shortLinkRepository.Add("sa42d45ds2", "https:\\\\testing.com\\ksjadkjas")
	if err != nil {
		t.Errorf("ShortLinkFixture: Error = '%v'", err.Error())
	}

	_, err = shortLinkRepository.Add("qwetyr123iu", "https:\\\\google.com")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	handler := NewGetShortLinkHandler(service.NewReadShortLinkService(shortLinkRepository))

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
			name:    "SuccessGettingShortLink_1",
			method:  http.MethodGet,
			request: "/sa42d45ds2",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"https:\\\\testing.com\\ksjadkjas\">Temporary Redirect</a>.\n\n",
				location:   "https:\\\\testing.com\\ksjadkjas",
			},
		},
		{
			name:    "SuccessGettingShortLink_2",
			method:  http.MethodGet,
			request: "/qwetyr123iu",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"https:\\\\google.com\">Temporary Redirect</a>.\n\n",
				location:   "https:\\\\google.com",
			},
		},
		{
			name:    "UsageMethodNotAllowed",
			method:  http.MethodPost,
			request: "/asd",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "405 method not allowed",
				location:   "",
			},
		},
		{
			name:    "UndefinedURL",
			method:  http.MethodGet,
			request: "/testesttest",
			body:    "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "not found by short link: testesttest",
				location:   "",
			},
		},
		{
			name:    "UncorrectURL",
			method:  http.MethodGet,
			request: "/",
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
			context.AddParam("id", strings.TrimLeft(tt.request, "/"))
			context.Request = request

			r := setupRouter()
			r.GET(`/:id`, handler.Handle)
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

			if location := tt.want.location; location != response.Header().Get("Location") {
				t.Errorf("%v: Content = '%v', want = '%v'", tt.name, location, response.Header().Get("Location"))
			}
		})
	}
}
