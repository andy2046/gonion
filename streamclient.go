package gonion

import (
	"context"

	"google.golang.org/grpc"
)

// NewStreamClient creates a new StreamClient chain.
func NewStreamClient(interceptors ...grpc.StreamClientInterceptor) StreamClientChain {
	return StreamClientChain{append(([]grpc.StreamClientInterceptor)(nil), interceptors...)}
}

// Then chains middlewares and returns the final grpc.StreamClientInterceptor.
func (scc StreamClientChain) Then(ssi grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	if ssi == nil {
		ssi = func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return streamer(ctx, desc, cc, method, opts...)
		}
	}
	newInterceptors := make([]grpc.StreamClientInterceptor, 0, len(scc.interceptors)+1)
	newInterceptors = append(newInterceptors, scc.interceptors...)
	newInterceptors = append(newInterceptors, ssi)
	newscc := StreamClientChain{newInterceptors}

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return newscc.chain(ctx, desc, cc, method, streamer, 0, opts...)
	}
}

func (scc StreamClientChain) chain(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, i int, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if i == len(scc.interceptors) {
		return streamer(ctx, desc, cc, method, opts...)
	}
	return scc.interceptors[i](ctx, desc, cc, method, func(ctx2 context.Context, desc2 *grpc.StreamDesc,
		cc2 *grpc.ClientConn, method2 string, opts2 ...grpc.CallOption) (grpc.ClientStream, error) {
		return scc.chain(ctx2, desc2, cc2, method2, streamer, i+1, opts2...)
	}, opts...)
}

// Append extends a chain, by appending interceptors to the request flow.
func (scc StreamClientChain) Append(interceptors ...grpc.StreamClientInterceptor) StreamClientChain {
	ins := make([]grpc.StreamClientInterceptor, 0, len(scc.interceptors)+len(interceptors))
	ins = append(ins, scc.interceptors...)
	ins = append(ins, interceptors...)

	return StreamClientChain{interceptors: ins}
}

// Extend extends a chain, by appending another chain to the request flow.
func (scc StreamClientChain) Extend(scchain StreamClientChain) StreamClientChain {
	return scc.Append(scchain.interceptors...)
}
