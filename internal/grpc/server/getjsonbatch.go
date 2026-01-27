package server

import (
	"context"

	"github.com/ArtShib/urlshortener/gen/go/shortener"
	"github.com/ArtShib/urlshortener/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverAPI) ListUserURLs(ctx context.Context, empty *emptypb.Empty) (*shortener.UserURLsResponse, error) {

	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.InvalidArgument, "userID empty")
	}

	urlsBatch, err := s.service.GetJSONBatch(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if len(urlsBatch) == 0 {
		return nil, status.Error(codes.InvalidArgument, "your short URLs not found")
	}
	urlData := make([]*shortener.URLData, len(urlsBatch))
	for i, url := range urlsBatch {
		urlData[i] = &shortener.URLData{}
		urlData[i].SetShortUrl(url.ShortURL)
		urlData[i].SetOriginalUrl(url.OriginalURL)
	}
	var resp shortener.UserURLsResponse
	resp.SetUrl(urlData)
	return &resp, nil
}
