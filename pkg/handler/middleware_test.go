package handler

import (
	"bookshelf-api/pkg/lib/slogdiscard"
	"bookshelf-api/pkg/service"
	"bookshelf-api/pkg/service/mocks"
	"errors"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehaviour func(auth *mocks.Authorization, token string)

	tests := []struct {
		name           string
		headerName     string
		headerValue    string
		token          string
		mockBehaviour  mockBehaviour
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(auth *mocks.Authorization, token string) {
				auth.
					On("ParseToken", token).
					Return(1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "1",
		},
		{
			name:           "invalid header name",
			headerName:     "",
			headerValue:    "Bearer token",
			token:          "token",
			mockBehaviour:  func(auth *mocks.Authorization, token string) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"error\":\"empty auth header\"}\n",
		},
		{
			name:           "invalid header value",
			headerName:     "Authorization",
			headerValue:    "Bear token",
			token:          "token",
			mockBehaviour:  func(auth *mocks.Authorization, token string) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"error\":\"invalid auth header\"}\n",
		},
		{
			name:           "empty token",
			headerName:     "Authorization",
			headerValue:    "Bearer ",
			token:          "token",
			mockBehaviour:  func(auth *mocks.Authorization, token string) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"error\":\"invalid auth header\"}\n",
		},
		{
			name:        "parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(auth *mocks.Authorization, token string) {
				auth.On("ParseToken", token).
					Return(0, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"error\":\"invalid token\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := mocks.NewAuthorization(t)
			tt.mockBehaviour(auth, tt.token)
			services := &service.Service{Authorization: auth}
			h := Handler{services}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := r.Context().Value("userID").(int)
				if !ok {
					t.Error("user id not found")
					render.Status(r, http.StatusInternalServerError)
					render.JSON(w, r, Error("user id not found"))
					return
				}
				render.Status(r, http.StatusOK)
				render.Data(w, r, []byte(strconv.Itoa(userID)))
			})
			handlerToTest := h.userIdentity(slogdiscard.NewDiscardLogger())(nextHandler)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(tt.headerName, tt.headerValue)
			handlerToTest.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
