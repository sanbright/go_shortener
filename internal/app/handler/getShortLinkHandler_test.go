package handler

import (
	"io"
	"sanbright/go_shortener/internal/app/middleware"
	"strings"
	"testing"

	"go.uber.org/zap"

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

func setupLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}

func TestGetShortLinkHandler_Handle(t *testing.T) {

	shortLinkRepository := repository.NewShortLinkRepository()

	_, err := shortLinkRepository.Add("sa42d45ds2", "https:\\\\testing.com\\ksjadkjas", "4c1b4334-8f1c-4874-8750-c5214e2f48b9")
	if err != nil {
		t.Errorf("ShortLinkFixture: Error = '%v'", err.Error())
	}

	_, err = shortLinkRepository.Add("qwetyr123iu", "https:\\\\google.com", "4c1b4334-8f1c-4874-8750-c5214e2f48b9")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	_, err = shortLinkRepository.Add("qwetyr123i1", "https:\\\\google1.com", "4c1b4334-8f1c-4874-8750-c5214e2f48b9")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	del := []string{"qwetyr123i1"}
	err = shortLinkRepository.Delete(del, "4c1b4334-8f1c-4874-8750-c5214e2f48b9")
	if err != nil {
		t.Errorf("hortLinkFixture: Error = '%v'", err.Error())
	}

	handler := NewGetShortLinkHandler(service.NewReadShortLinkService(shortLinkRepository))
	logger := setupLogger()
	loggerMiddleware := middleware.Logger(logger)

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
			name:    "NotFoundGettingShortLink",
			method:  http.MethodGet,
			request: "/qwetyr123i1",
			want: want{
				statusCode: http.StatusGone,
				body:       "Not found link",
				location:   "",
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
			context.AddParam("userId", "4c1b4334-8f1c-4874-8750-c5214e2f48b9")
			context.Request = request

			r := setupRouter()

			r.GET(`/:id`, loggerMiddleware, handler.Handle)
			r.ServeHTTP(response, request)

			if code := tt.want.statusCode; code != response.Code {
				t.Errorf("%v: StatusCode = '%v', want = '%v'", tt.name, response.Code, code)
			}

			body, err := io.ReadAll(response.Body)

			if err != nil {
				t.Errorf("%v: Error = '%v'", tt.name, err.Error())
			}

			if tbody := tt.want.body; tbody != string(body) {
				t.Errorf("%v: Content_body = '%v', want = '%v'", tt.name, tbody, string(body))
			}

			if location := tt.want.location; location != response.Header().Get("Location") {
				t.Errorf("%v: Content = '%v', want = '%v'", tt.name, location, response.Header().Get("Location"))
			}
		})
	}
}
