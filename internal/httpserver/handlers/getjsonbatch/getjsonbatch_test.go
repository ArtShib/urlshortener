package getjsonbatch

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.URLUserBatch), args.Error(1)
}

func TestGetJSONBatchHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name                string
		userID              string
		mockFunc            func(m *MockURLService, userID string)
		expectedStatus      int
		expectedContentType string
		expectedBody        string
		isJSONResponse      bool
	}{
		{
			name:   "Success",
			userID: "2",
			mockFunc: func(m *MockURLService, userID string) {
				m.On("GetJSONBatch", mock.Anything, userID).
					Return(model.URLUserBatch{
						model.URLUser{ShortURL: "http://localhost/sdfdfg", OriginalURL: "https://google.com"},
						model.URLUser{ShortURL: "http://localhost/asdasd", OriginalURL: "https://yandex.ru"}}, nil).
					Once()
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedBody:        `[{"short_url":"http://localhost/sdfdfg","original_url":"https://google.com"},{"short_url":"http://localhost/asdasd","original_url":"https://yandex.ru"}]`,
			isJSONResponse:      true,
		},
		{
			name:   "NoContent)",
			userID: "2",
			mockFunc: func(m *MockURLService, userID string) {
				m.On("GetJSONBatch", mock.Anything, userID).
					Return(model.URLUserBatch{}, nil).
					Once()
			},
			expectedStatus:      http.StatusNoContent,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "No Content\n",
			isJSONResponse:      false,
		},
		{
			name:                "Unauthorized",
			userID:              "",
			mockFunc:            func(m *MockURLService, userID string) {},
			expectedStatus:      http.StatusUnauthorized,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "Unauthorized\n",
			isJSONResponse:      false,
		},
		{
			name:   "InternalError",
			userID: "2",
			mockFunc: func(m *MockURLService, userID string) {
				m.On("GetJSONBatch", mock.Anything, userID).
					Return(model.URLUserBatch{}, errors.New(http.StatusText(http.StatusInternalServerError))).
					Once()
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBody:        "Internal Server Error\n",
			isJSONResponse:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := new(MockURLService)

			test.mockFunc(svc, test.userID)

			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

			if test.userID != "" {
				ctx := context.WithValue(req.Context(), model.UserIDKey, test.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()

			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, test.expectedContentType, resp.Header.Get("Content-Type"))
			if test.isJSONResponse {
				assert.JSONEq(t, test.expectedBody, string(resBody))
			} else {
				assert.Equal(t, test.expectedBody, string(resBody))
			}

			svc.AssertExpectations(t)
		})
	}
}
