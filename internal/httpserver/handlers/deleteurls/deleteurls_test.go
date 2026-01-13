package deleteurls

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) AddRequest(req model.DeleteRequest) {
	m.Called(req)
}

func TestGetJSONBatchHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		inputBody      string
		userID         string
		mockFunc       func(m *MockURLService)
		expectedStatus int
	}{
		{
			name:      "Success",
			inputBody: `["aedsadd","restfgt"]`,
			userID:    "2",
			mockFunc: func(m *MockURLService) {
				m.On("AddRequest", mock.MatchedBy(func(req model.DeleteRequest) bool {
					return req.UserID == "2" && len(req.UUIDs) == 2 && req.UUIDs[0] == "aedsadd"
				})).
					Once()
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Unauthorized",
			inputBody:      `["aedsadd","restfgt"]`,
			userID:         "",
			mockFunc:       func(m *MockURLService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "InternalError",
			inputBody:      `["aedsadd]`,
			userID:         "2",
			mockFunc:       func(m *MockURLService) {},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := new(MockURLService)

			test.mockFunc(svc)

			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString(test.inputBody))

			if test.userID != "" {
				ctx := context.WithValue(req.Context(), model.UserIDKey, test.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()

			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			svc.AssertExpectations(t)
		})
	}
}
