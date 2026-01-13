package ping

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestPingHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		mockFunc       func(m *MockURLService)
		expectedStatus int
	}{
		{
			name: "Success",
			mockFunc: func(m *MockURLService) {
				m.On("Ping", mock.Anything).
					Return(nil).
					Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "InternalError",
			mockFunc: func(m *MockURLService) {
				m.On("Ping", mock.Anything).
					Return(errors.New(http.StatusText(http.StatusInternalServerError))).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := new(MockURLService)

			test.mockFunc(svc)
			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, test.expectedStatus, resp.StatusCode)
			svc.AssertExpectations(t)
		})
	}
}
