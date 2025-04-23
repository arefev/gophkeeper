package interceptor

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (m *middleware) StreamCheckToken(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var token string
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing context")
	}

	values := md.Get("token")
	if len(values) > 0 {
		token = values[0]
	}

	if token == "" {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	m.app.Log.Debug("incoming context", zap.String("token", token))
	err := handler(srv, ss)
	return err
}
