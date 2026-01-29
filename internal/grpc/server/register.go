package server

import (
	"google.golang.org/grpc"

	shortener "github.com/ArtShib/urlshortener/gen/go/shortener"
)

func Register(gRPC *grpc.Server, svc URLService) {
	shortener.RegisterShortenerServiceServer(gRPC, &serverAPI{
		service: svc,
	})
}
