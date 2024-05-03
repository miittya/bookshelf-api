package handler

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/lib/slogdiscard"
	"bookshelf-api/pkg/service"
	"bookshelf-api/pkg/service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateList(t *testing.T) {
	type mockBehaviour func(list *mocks.List, userID int, input bookshelf.List)

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		inputBody      string
		inputList      bookshelf.List
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "OK",
			mockBehaviour: func(list *mocks.List, userID int, input bookshelf.List) {
				list.On("Create", userID, input).Return(1, nil)
			},
			inputBody: `{"title":"title","description":"description"}`,
			inputList: bookshelf.List{
				Title:       "title",
				Description: "description",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"list_id\":1}\n",
		},
		{
			name: "Only Title",
			mockBehaviour: func(list *mocks.List, userID int, input bookshelf.List) {
				list.On("Create", userID, input).Return(1, nil)
			},
			inputBody: `{"title":"title"}`,
			inputList: bookshelf.List{
				Title: "title",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"list_id\":1}\n",
		},
		{
			name:          "Only Description",
			mockBehaviour: func(list *mocks.List, userID int, input bookshelf.List) {},
			inputBody:     `{"description":"description"}`,
			inputList: bookshelf.List{
				Description: "description",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error\":\"invalid request\"}\n",
		},
		{
			name: "service error",
			mockBehaviour: func(list *mocks.List, userID int, input bookshelf.List) {
				list.On("Create", userID, input).Return(0, errors.New("service error"))
			},
			inputBody: `{"title":"title","description":"description"}`,
			inputList: bookshelf.List{
				Title:       "title",
				Description: "description",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error\":\"cannot create list\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := mocks.NewList(t)
			tt.mockBehaviour(list, 1, tt.inputList)
			services := &service.Service{List: list}
			handler := Handler{services}

			r := chi.NewRouter()
			r.Post("/", handler.createList(slogdiscard.NewDiscardLogger()))

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.inputBody)))
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), "userID", 1)
			r.ServeHTTP(w, req.WithContext(ctx))

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
