package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	mockService "acsp/internal/service/mocks"

	"acsp/internal/logs"
	"acsp/internal/model"
	"acsp/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(r *mockService.MockAuthorization, user model.User)

	tests := []struct {
		name                 string
		inputBody            map[string]any
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: map[string]any{
				"email":    "eldos2020@gmail.com",
				"name":     "testname",
				"password": "qwerty123",
			},
			inputUser: model.User{
				Email:    "eldos2020@gmail.com",
				Name:     "testname",
				Password: "qwerty123",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":"1"}`,
		},
		{
			name: "Service Error",
			inputBody: map[string]any{
				"email":    "eldos2020@gmail.com",
				"name":     "testname",
				"password": "qwerty123",
			},
			inputUser: model.User{
				Email:    "eldos2020@gmail.com",
				Name:     "testname",
				Password: "qwerty123",
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name: "Wrong Input",
			inputBody: map[string]any{
				"email": "eldos2020@gmail.com",
			},
			inputUser:            model.User{},
			mockBehavior:         func(r *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			_ = logs.InitLogger()
			repo := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(repo, testCase.inputUser)

			services := &service.Service{Authorization: repo}
			handler := Handler{services}

			app := fiber.New()
			app.Post("/sign-up", handler.signUp)

			requestBody, _ := json.Marshal(testCase.inputBody)

			request := httptest.NewRequest("POST", "/sign-up",
				bytes.NewReader(requestBody))
			request.Header.Add("Content-Type", "application/json")

			response, _ := app.Test(request)
			data, _ := io.ReadAll(response.Body)

			assert.Equal(t, response.StatusCode, testCase.expectedStatusCode)
			assert.Matches(t, string(data), testCase.expectedResponseBody)
		})
	}
}
