package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"acsp/internal/service"
	mockService "acsp/internal/service/mocks"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(r *mockService.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Users",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mockService.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return("1", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Users",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Users",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Users",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mockService.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return("0", errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAuthorization(c)
			test.mockBehavior(repo, test.token)

			services := &service.Service{Authorization: repo}
			handler := Handler{services}

			app := fiber.New()
			app.Get("/identity", handler.userIdentity, func(ctx *fiber.Ctx) error {
				id := ctx.GetRespHeader(userCtx)
				return ctx.Status(http.StatusOK).SendString(id)
			})

			request := httptest.NewRequest("GET", "/identity", nil)
			request.Header.Add("Content-Type", "application/json")
			request.Header.Set(test.headerName, test.headerValue)

			response, _ := app.Test(request)
			data, _ := io.ReadAll(response.Body)

			assert.Equal(t, response.StatusCode, test.expectedStatusCode)
			assert.Equal(t, string(data), test.expectedResponseBody)
		})
	}
}
