package server

import (
	"context"

	"github.com/ArtShib/urlshortener/gen/go/shortener"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) ExpandURL(ctx context.Context, req *shortener.URLExpandRequest) (*shortener.URLExpandResponse, error) {
	shortCode := req.GetId()
	if shortCode == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	url, err := s.service.GetID(ctx, shortCode)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var res shortener.URLExpandResponse
	res.SetResult(url.OriginalURL)
	return &res, nil
}
