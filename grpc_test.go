package gonion

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	contextKey string

	testServerStream struct {
		grpc.ServerStream
		ctx context.Context
	}

	testClientStream struct {
		grpc.ClientStream
	}
)

var (
	someServiceName = "Service.Method"
	someValue       = 1
	parentContext   = context.WithValue(context.Background(), contextKey("parent"), someValue)
)

func (ss testServerStream) Context() context.Context {
	return ss.ctx
}

func (testServerStream) SendMsg(m interface{}) error {
	return nil
}

func (testServerStream) RecvMsg(m interface{}) error {
	return nil
}

func requireContextValue(ctx context.Context, t *testing.T, key string, msg ...interface{}) {
	val := ctx.Value(contextKey(key))
	require.NotNil(t, val, msg...)
	require.Equal(t, someValue, val, msg...)
}

func TestUnaryServer(t *testing.T) {
	input := "input"

	foo := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if got, want := req.(string), input; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		ctx = context.WithValue(ctx, contextKey("foo"), 1)
		return handler(ctx, req)
	}
	bar := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if got, want := req.(string), input; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value(contextKey("foo")).(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		ctx = context.WithValue(ctx, contextKey("bar"), 2)
		return handler(ctx, req)
	}

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "Test.UnaryServer",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if got, want := req.(string), input; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value(contextKey("foo")).(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value(contextKey("bar")).(int), 2; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		return "output", nil
	}

	ctx := context.Background()
	interceptor := NewUnaryServer(foo, bar).Then(nil)
	out, err := interceptor(ctx, input, unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out, "output"; got != want {
		t.Errorf("expect output to %q, but got %q", want, got)
	}
}

func TestStreamServer(t *testing.T) {
	foo := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := stream.Context()
		ctx = context.WithValue(ctx, contextKey("foo"), 1)
		stream = testServerStream{stream, ctx}
		return handler(srv, stream)
	}
	bar := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := stream.Context()
		if got, want := ctx.Value(contextKey("foo")).(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}

		ctx = context.WithValue(ctx, contextKey("bar"), 2)
		stream = testServerStream{stream, ctx}
		return handler(srv, stream)
	}

	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "Test.StreamServer",
		IsServerStream: true,
	}
	streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
		ctx := stream.Context()
		if got, want := ctx.Value(contextKey("foo")).(int), 1; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		if got, want := ctx.Value(contextKey("bar")).(int), 2; got != want {
			t.Errorf("expect output to %q, but got %q", want, got)
		}
		return nil
	}
	testService := struct{}{}
	testStream := &testServerStream{ctx: context.Background()}

	interceptor := NewStreamServer(foo, bar).Then(nil)
	err := interceptor(testService, testStream, streamInfo, streamHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnaryClient(t *testing.T) {
	ignoredMd := metadata.Pairs("foo", "bar")
	parentOpts := []grpc.CallOption{grpc.Header(&ignoredMd)}
	reqMessage := "request"
	replyMessage := "reply"
	outputError := fmt.Errorf("some error")

	first := func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requireContextValue(ctx, t, "parent", "first must know the parent context value")
		require.Equal(t, someServiceName, method, "first must know someService")
		require.Len(t, opts, 1, "first should see parent CallOptions")
		wrappedCtx := context.WithValue(ctx, contextKey("first"), 1)
		return invoker(wrappedCtx, method, req, reply, cc, opts...)
	}
	second := func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requireContextValue(ctx, t, "parent", "second must know the parent context value")
		require.Equal(t, someServiceName, method, "second must know someService")
		require.Len(t, opts, 1, "second should see parent CallOptions")
		wrappedOpts := append(opts, grpc.FailFast(true))
		wrappedCtx := context.WithValue(ctx, contextKey("second"), 1)
		return invoker(wrappedCtx, method, req, reply, cc, wrappedOpts...)
	}
	invoker := func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		require.Equal(t, someServiceName, method, "invoker must know someService")
		requireContextValue(ctx, t, "parent", "invoker must know the parent context value")
		requireContextValue(ctx, t, "first", "invoker must know the first context value")
		requireContextValue(ctx, t, "second", "invoker must know the second context value")
		require.Len(t, opts, 2, "invoker should see both CallOpts from second and parent")
		return outputError
	}
	chain := NewUnaryClient(first, second).Then(nil)
	err := chain(parentContext, someServiceName, reqMessage, replyMessage, nil, invoker, parentOpts...)
	require.Equal(t, outputError, err, "chain must return invokers's error")
}

func TestStreamClient(t *testing.T) {
	ignoredMd := metadata.Pairs("foo", "bar")
	parentOpts := []grpc.CallOption{grpc.Header(&ignoredMd)}
	clientStream := &testClientStream{}
	testStreamDesc := &grpc.StreamDesc{ClientStreams: true, ServerStreams: true, StreamName: someServiceName}

	first := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requireContextValue(ctx, t, "parent", "first must know the parent context value")
		require.Equal(t, someServiceName, method, "first must know someService")
		require.Len(t, opts, 1, "first should see parent CallOptions")
		wrappedCtx := context.WithValue(ctx, contextKey("first"), 1)
		return streamer(wrappedCtx, desc, cc, method, opts...)
	}
	second := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requireContextValue(ctx, t, "parent", "second must know the parent context value")
		require.Equal(t, someServiceName, method, "second must know someService")
		require.Len(t, opts, 1, "second should see parent CallOptions")
		wrappedOpts := append(opts, grpc.FailFast(true))
		wrappedCtx := context.WithValue(ctx, contextKey("second"), 1)
		return streamer(wrappedCtx, desc, cc, method, wrappedOpts...)
	}
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		require.Equal(t, someServiceName, method, "streamer must know someService")
		require.Equal(t, testStreamDesc, desc, "streamer must see the right StreamDesc")

		requireContextValue(ctx, t, "parent", "streamer must know the parent context value")
		requireContextValue(ctx, t, "first", "streamer must know the first context value")
		requireContextValue(ctx, t, "second", "streamer must know the second context value")
		require.Len(t, opts, 2, "streamer should see both CallOpts from second and parent")
		return clientStream, nil
	}
	chain := NewStreamClient(first, second).Then(nil)
	someStream, err := chain(parentContext, testStreamDesc, nil, someServiceName, streamer, parentOpts...)
	require.NoError(t, err, "chain must not return an error as nothing there reutrned it")
	require.Equal(t, clientStream, someStream, "chain must return invokers's clientstream")
}
