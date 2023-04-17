package interceptor

import (
	"context"
	"database/sql"
	"errors"
	"sample-grpc-server/server"

	"sample-grpc-server/database"

	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var ErrNoAccessToken = errors.New("session: access_token not found")

func RecoveryFunc(p interface{}) error {
	return status.Errorf(codes.Unknown, "unexpected error: %v", p)
}

func AuthInterceptor(db database.Querier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isAuthFree(info.FullMethod) {
			return handler(ctx, req)
		}

		newCtx, err := authenticate(ctx, db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Error(codes.Unauthenticated, "failed to authenticate")
			} else if errors.Is(err, ErrNoAccessToken) {
				return nil, status.Error(codes.Unauthenticated, "failed to authenticate")
			}
			return nil, status.Error(codes.Internal, "server error")
		}

		return handler(newCtx, req)
	}
}

func isAuthFree(method string) bool {
	authFreeMethods := []string{
		"/backend.BackendService/HelloWorld",
		"/backend.BackendService/SignUp",
		"/backend.BackendService/Login",
	}

	for _, m := range authFreeMethods {
		if m == method {
			return true
		}
	}

	return false
}

func authenticate(ctx context.Context, db database.Querier) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, xerrors.Errorf("failed to extract metadata")
	}

	tokens := md.Get("access_token")
	if len(tokens) < 1 {
		return ctx, ErrNoAccessToken
	}

	session, err := db.GetSession(ctx, database.GetSessionParams{AccessToken: tokens[0]})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.Errorf("session not found: %w", err)
		}
		return nil, xerrors.Errorf("failed to get session: %v", err)
	}

	ctx = context.WithValue(ctx, server.KeyUserID, session.ID)

	return ctx, nil
}
