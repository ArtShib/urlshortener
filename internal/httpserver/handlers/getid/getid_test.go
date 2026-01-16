package getid

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) GetID(ctx context.Context, shortCode string) (*model.URL, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.URL), args.Error(1)
}

func withURLParam(r *http.Request, key, value string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestGetIDHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name             string
		urlParamID       string
		mockFunc         func(m *MockURLService, shortCode string)
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:       "Success",
			urlParamID: "sdsd34vcx",
			mockFunc: func(m *MockURLService, shortCode string) {
				m.On("GetID", mock.Anything, shortCode).
					Return(&model.URL{OriginalURL: "https://google.com", DeletedFlag: false}, nil).
					Once()
			},
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://google.com",
		},
		{
			name:       "Gone",
			urlParamID: "sdsd34vcx",
			mockFunc: func(m *MockURLService, shortCode string) {
				m.On("GetID", mock.Anything, shortCode).
					Return(&model.URL{OriginalURL: "https://google.com", DeletedFlag: true}, nil).
					Once()
			},
			expectedStatus:   http.StatusGone,
			expectedLocation: "",
		},
		{
			name:             "EmptyID",
			urlParamID:       "",
			mockFunc:         func(m *MockURLService, shortCode string) {},
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
		},
		{
			name:       "InternalError",
			urlParamID: "sdsd34vcx",
			mockFunc: func(m *MockURLService, shortCode string) {
				m.On("GetID", mock.Anything, shortCode).
					Return(nil, errors.New("database error")).
					Once()
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedLocation: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := new(MockURLService)
			test.mockFunc(svc, test.urlParamID)

			handler := New(logger, svc)

			req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
			req = withURLParam(req, "shortCode", test.urlParamID)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer func() {
				if err := resp.Body.Close(); err != nil {
					require.NoError(t, err)
				}
			}()

			assert.Equal(t, test.expectedStatus, resp.StatusCode)
			if test.expectedLocation != "" {
				assert.Equal(t, test.expectedLocation, resp.Header.Get("Location"))
			} else {
				assert.Empty(t, resp.Header.Get("Location"))
			}
			svc.AssertExpectations(t)
		})
	}
}
