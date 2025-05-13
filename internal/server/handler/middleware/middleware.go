package middleware

import (
	"context"

	"google.golang.org/grpc"
)

// WrappedServerStream struct.
//
//	WrappedContext - context using in steaming.
type WrappedServerStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

// Context return context, wich was wrapped in steam interceptor.
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream create pointer on WrappedServerStream struct.
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
