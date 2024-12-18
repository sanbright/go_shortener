package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"sanbright/go_shortener/internal/app/generator"
	"testing"
)

func TestAuth(t *testing.T) {
	type want struct {
		statusCode int
		userId     string
	}

	tests := []struct {
		name       string
		authCookie string
		want       want
	}{
		{
			name:       "Unauthorized",
			authCookie: "I8LumVeMYJlq8p",
			want: want{
				statusCode: -1,
				userId:     "",
			},
		},
		{
			name:       "Empty_Cookie",
			authCookie: "",
			want: want{
				statusCode: -1,
				userId:     "",
			},
		},
		{
			name:       "Authorized",
			authCookie: "I8LumVeMYJlq8pNoeeY0s1EzbMS90OFaFnH0uXYKv3I7FEbDBSDPMvRjLDgVZx3Q8wGSGA==",
			want: want{
				statusCode: http.StatusOK,
				userId:     "653e7307-6960-4b60-ab1b-44cd2f662634",
			},
		},
	}

	log := setupLogger()
	cGen := generator.NewCryptGenerator("$$ecuRityKe453H@")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := Auth(cGen, log)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authCookie != "" {
				request.AddCookie(&http.Cookie{Name: "Auth", Value: tt.authCookie})
			}

			request.AddCookie(&http.Cookie{Name: "Some", Value: "Value"})

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = request
			m(ctx)

			if tt.want.statusCode != -1 && ctx.Writer.Status() != tt.want.statusCode {
				t.Errorf("status = %v, want %v", ctx.Writer.Status(), tt.want.statusCode)
			}

			userID, _ := ctx.Get("UserID")

			if userID != nil && userID != tt.want.userId {
				t.Errorf("userID = %v, want %v", userID, tt.want.userId)
			}

		})
	}
}
