package interceptors

import (
	"context"

	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(auth *auth.Service) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		var userID string
		needAuth := true
		var err error
		mb, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := mb.Get("authorization")
			if len(values) != 0 {
				if auth.ValidateToken(values[0]) {
					userID = auth.GetUserID(values[0])
					needAuth = false
				}
			}
		}
		if needAuth {
			userID, err = auth.GenerateUserID()
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "failed to generate user id")
			}
		}

		ctx = context.WithValue(ctx, model.UserIDKey, userID)
		return handler(ctx, req)
	}
}
