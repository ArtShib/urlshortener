package server

import (
	"context"

	"github.com/ArtShib/urlshortener/gen/go/shortener"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) ShortenURL(ctx context.Context, req *shortener.URLShortenRequest) (*shortener.URLShortenResponse, error) {
	url := req.GetUrl()

	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is empty")
	}
	shortURL, err := s.service.Shorten(ctx, url)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var resp shortener.URLShortenResponse
	resp.SetResult(shortURL)
	return &resp, nil

}
