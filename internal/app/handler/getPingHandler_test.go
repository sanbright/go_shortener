package handler

import (
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/config"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPingHandler_Handle(t *testing.T) {
	configuration, _ := config.NewConfig("localhost:8080", "", "", "")
	handler := NewGetPingHandler(configuration)

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
			name:    "SuccessGettingShortLink_2",
			method:  http.MethodGet,
			request: "/ping",
			want: want{
				statusCode: http.StatusBadGateway,
				body:       "<a href=\"https:\\\\google.com\">Temporary Redirect</a>.\n\n",
				location:   "https:\\\\google.com",
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
			r.GET(`/ping`, handler.Handle)
			//r.ServeHTTP(response, request)

			//if code := tt.want.statusCode; code != response.Code {
			//	t.Errorf("%v: StatusCode = '%v', want = '%v'", tt.name, response.Code, code)
			//}

			//body, err := io.ReadAll(response.Body)

			//if err != nil {
			//	t.Errorf("%v: Error = '%v'", tt.name, err.Error())
			//}

			//if tbody := tt.want.body; tbody != string(body) {
			//	t.Errorf("%v: Content = '%v', want = '%v'", tt.name, tbody, string(body))
			//}
			//
			//if location := tt.want.location; location != response.Header().Get("Location") {
			//	t.Errorf("%v: Content = '%v', want = '%v'", tt.name, location, response.Header().Get("Location"))
			//}
		})
	}
}