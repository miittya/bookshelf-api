package handler

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/lib/slogdiscard"
	"bookshelf-api/pkg/service"
	"bookshelf-api/pkg/service/mocks"
	"bytes"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehaviour func(auth *mocks.Authorization, user bookshelf.User)

	tests := []struct {
		name           string
		inputBody      string
		inputUser      bookshelf.User
		mockBehaviour  mockBehaviour
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test","password":"qwerty"}`,
			inputUser: bookshelf.User{
				Username: "test",
				Password: "qwerty",
			},
			mockBehaviour: func(auth *mocks.Authorization, user bookshelf.User) {
				auth.
					On("CreateUser", user).
					Return(1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1}\n",
		},
		{
			name:      "No password",
			inputBody: `{"username":"test"}`,
			inputUser: bookshelf.User{
				Username: "test",
			},
			mockBehaviour:  func(auth *mocks.Authorization, user bookshelf.User) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error\":\"invalid request\"}\n",
		},
		{
			name:      "No username",
			inputBody: `{"password":"qwerty"}`,
			inputUser: bookshelf.User{
				Password: "qwerty",
			},
			mockBehaviour:  func(auth *mocks.Authorization, user bookshelf.User) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error\":\"invalid request\"}\n",
		},
		{
			name:           "Empty request",
			inputBody:      "{}",
			inputUser:      bookshelf.User{},
			mockBehaviour:  func(auth *mocks.Authorization, user bookshelf.User) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error\":\"invalid request\"}\n",
		},
		{
			name:      "Service error",
			inputBody: `{"username":"test","password":"qwerty"}`,
			inputUser: bookshelf.User{
				Username: "test",
				Password: "qwerty",
			},
			mockBehaviour: func(auth *mocks.Authorization, user bookshelf.User) {
				auth.
					On("CreateUser", user).
					Return(0, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error\":\"cannot create user\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := mocks.NewAuthorization(t)
			tt.mockBehaviour(auth, tt.inputUser)
			services := &service.Service{Authorization: auth}
			h := Handler{services}

			r := chi.NewRouter()
			r.Post("/sign-up", h.SignUp(slogdiscard.NewDiscardLogger()))

			req, _ := http.NewRequest(http.MethodPost, "/sign-up", bytes.NewReader([]byte(tt.inputBody)))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
