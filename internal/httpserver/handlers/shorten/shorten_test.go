package shorten

import (
	"bytes"
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

func (m *MockURLService) Shorten(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func TestShortenHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		inputBody      string
		mockFunc       func(m *MockURLService, body string)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success",
			inputBody: "https://google.com",
			mockFunc: func(m *MockURLService, body string) {
				m.On("Shorten", mock.Anything, body).
					Return("http://localhost/sdfdfg", nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost/sdfdfg",
		},
		{
			name:      "Conflict",
			inputBody: "https://google.com",
			mockFunc: func(m *MockURLService, body string) {
				m.On("Shorten", mock.Anything, body).
					Return("http://localhost/sdfdfg", model.ErrURLConflict).
					Once()
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   "http://localhost/sdfdfg",
		},
		{
			name:      "InternalError",
			inputBody: "https://google.com",
			mockFunc: func(m *MockURLService, body string) {
				m.On("Shorten", mock.Anything, body).
					Return("", errors.New(http.StatusText(http.StatusInternalServerError))).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := new(MockURLService)
			test.mockFunc(svc, test.inputBody)

			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.inputBody))
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			defer func() {
				if err := resp.Body.Close(); err != nil {
					require.NoError(t, err)
				}
			}()

			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.Equal(t, test.expectedBody, string(resBody))

			svc.AssertExpectations(t)
		})
	}
}
