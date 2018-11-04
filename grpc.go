package gonion

import (
	"context"

	"google.golang.org/grpc"
)

type (
	// StreamClientChain chains grpc StreamClient middlewares.
	StreamClientChain struct {
		interceptors []grpc.StreamClientInterceptor
	}

	// StreamServerChain chains grpc StreamServer middlewares.
	StreamServerChain struct {
		interceptors []grpc.StreamServerInterceptor
	}

	// UnaryClientChain chains grpc UnaryClient middlewares.
	UnaryClientChain struct {
		interceptors []grpc.UnaryClientInterceptor
	}

	// UnaryServerChain chains grpc UnaryServer middlewares.
	UnaryServerChain struct {
		interceptors []grpc.UnaryServerInterceptor
	}

	// WrappedServerStream is a wrapped grpc.ServerStream with context.
	WrappedServerStream struct {
		grpc.ServerStream
		WrappedContext context.Context
	}
)

// WithUnaryServer returns a grpc.Server config option with multiple unary interceptors.
func WithUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.UnaryInterceptor(NewUnaryServer(interceptors...).Then(nil))
}

// WithStreamServer returns a grpc.Server config option with multiple stream interceptors.
func WithStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.StreamInterceptor(NewStreamServer(interceptors...).Then(nil))
}

// Context returns the wrapper's WrappedContext.
func (w WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream with the ability to overwrite context.
func WrapServerStream(stream grpc.ServerStream) WrappedServerStream {
	if wss, ok := stream.(WrappedServerStream); ok {
		return wss
	}
	return WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
