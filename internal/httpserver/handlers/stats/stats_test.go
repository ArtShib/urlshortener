package stats

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
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Stats(ctx context.Context) (*model.Stats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Stats), args.Error(1)
}

func TestNew(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		mockSetup      func(m *MockURLService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			mockSetup: func(m *MockURLService) {
				m.On("Stats", mock.Anything).Return(&model.Stats{
					CountURLs:  100,
					CountUsers: 50,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"urls":100,"users":50}`,
		},
		{
			name: "Service Error",
			mockSetup: func(m *MockURLService) {
				m.On("Stats", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   http.StatusText(http.StatusInternalServerError),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSvc := new(MockURLService)
			test.mockSetup(mockSvc)

			handler := New(logger, mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)

			ctx := context.WithValue(req.Context(), model.UserIDKey, "sad213d")
			req = req.WithContext(ctx)

			r := httptest.NewRecorder()

			handler.ServeHTTP(r, req)

			assert.Equal(t, test.expectedStatus, r.Code)

			if test.expectedStatus == http.StatusOK {
				assert.JSONEq(t, test.expectedBody, r.Body.String())
			} else {
				assert.Contains(t, r.Body.String(), test.expectedBody)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}
