package shortenjsonbatch

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

func (m *MockURLService) ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error) {
	args := m.Called(ctx, urls)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.ResponseShortenerBatchArray), args.Error(1)
}

func TestShortenJSONHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		inputBody      string
		mockFunc       func(m *MockURLService, urls model.RequestShortenerBatchArray)
		expectedStatus int
		expectedBody   string
		isJSONResponse bool
	}{
		{
			name:      "Success",
			inputBody: `[{"correlation_id":"uuid1","original_url":"https://google.com"},{"correlation_id":"uuid2","original_url":"https://yandex.ru"}]`,
			mockFunc: func(m *MockURLService, urls model.RequestShortenerBatchArray) {
				m.On("ShortenJSONBatch", mock.Anything, urls).
					Return(model.ResponseShortenerBatchArray{
						model.ResponseShortenerBatch{CorrelationID: "uuid1", ShortURL: "http://localhost/sdfdfg"},
						model.ResponseShortenerBatch{CorrelationID: "uuid2", ShortURL: "http://localhost/asdasd"}}, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `[{"correlation_id":"uuid1","short_url":"http://localhost/sdfdfg"},{"correlation_id":"uuid2","short_url":"http://localhost/asdasd"}]`,
			isJSONResponse: true,
		},
		{
			name:      "InternalError",
			inputBody: `[{"correlation_id":"uuid1","original_url":"https://google.com"},{"correlation_id":"uuid2","original_url":"https://yandex.ru"}]`,
			mockFunc: func(m *MockURLService, urls model.RequestShortenerBatchArray) {
				m.On("ShortenJSONBatch", mock.Anything, urls).
					Return(model.ResponseShortenerBatchArray{}, errors.New(http.StatusText(http.StatusInternalServerError))).
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

			var reqBody model.RequestShortenerBatchArray
			if err := json.Unmarshal([]byte(test.inputBody), &reqBody); err != nil {
				require.NoError(t, err)
			}
			test.mockFunc(svc, reqBody)

			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(test.inputBody))
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
			defer func() {
				if err := resp.Body.Close(); err != nil {
					require.NoError(t, err)
				}
			}()

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
