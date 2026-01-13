package shortenjson

import (
	"bytes"
	"context"
	"encoding/json"
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

func (m *MockURLService) ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResponseShortener), args.Error(1)
}

func TestShortenJSONHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		inputBody      string
		mockFunc       func(m *MockURLService, body string)
		expectedStatus int
		expectedBody   string
		isJSONResponse bool
	}{
		{
			name:      "Success",
			inputBody: `{"url": "https://google.com"}`,
			mockFunc: func(m *MockURLService, body string) {
				m.On("ShortenJSON", mock.Anything, body).
					Return(&model.ResponseShortener{Result: "http://localhost/sdfdfg"}, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"result": "http://localhost/sdfdfg"}`,
			isJSONResponse: true,
		},
		{
			name:      "Conflict",
			inputBody: `{"url": "https://google.com"}`,
			mockFunc: func(m *MockURLService, body string) {
				m.On("ShortenJSON", mock.Anything, body).
					Return(&model.ResponseShortener{Result: "http://localhost/sdfdfg"}, model.ErrURLConflict).
					Once()
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"result": "http://localhost/sdfdfg"}`,
			isJSONResponse: true,
		},
		{
			name:      "InternalError",
			inputBody: `{"url": "https://google.com"}`,
			mockFunc: func(m *MockURLService, body string) {
				m.On("ShortenJSON", mock.Anything, body).
					Return(&model.ResponseShortener{}, errors.New(http.StatusText(http.StatusInternalServerError))).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error\n",
			isJSONResponse: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := new(MockURLService)

			var reqBody model.RequestShortener
			if err := json.Unmarshal([]byte(test.inputBody), &reqBody); err != nil {
				require.NoError(t, err)
			}
			test.mockFunc(svc, reqBody.URL)

			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(test.inputBody))
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
			if test.isJSONResponse {
				assert.JSONEq(t, test.expectedBody, string(resBody))
			} else {
				assert.Equal(t, test.expectedBody, string(resBody))
			}

			svc.AssertExpectations(t)
		})
	}
}
