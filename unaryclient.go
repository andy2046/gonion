package gonion

import (
	"context"

	"google.golang.org/grpc"
)

// NewUnaryClient creates a new UnaryClient chain.
func NewUnaryClient(interceptors ...grpc.UnaryClientInterceptor) UnaryClientChain {
	return UnaryClientChain{append(([]grpc.UnaryClientInterceptor)(nil), interceptors...)}
}

// Then chains middlewares and returns the final grpc.UnaryClientInterceptor.
func (ucc UnaryClientChain) Then(usi grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	if usi == nil {
		usi = func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
	newInterceptors := make([]grpc.UnaryClientInterceptor, 0, len(ucc.interceptors)+1)
	newInterceptors = append(newInterceptors, ucc.interceptors...)
	newInterceptors = append(newInterceptors, usi)
	newucc := UnaryClientChain{newInterceptors}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return newucc.chain(ctx, method, req, reply, cc, invoker, 0, opts...)
	}
}

func (ucc UnaryClientChain) chain(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, i int, opts ...grpc.CallOption) error {
	if i == len(ucc.interceptors) {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	return ucc.interceptors[i](ctx, method, req, reply, cc, func(ctx2 context.Context, method2 string, req2, reply2 interface{},
		cc2 *grpc.ClientConn, opts2 ...grpc.CallOption) error {
		return ucc.chain(ctx2, method2, req2, reply2, cc2, invoker, i+1, opts2...)
	}, opts...)
}

// Append extends a chain, by appending interceptors to the request flow.
func (ucc UnaryClientChain) Append(interceptors ...grpc.UnaryClientInterceptor) UnaryClientChain {
	ins := make([]grpc.UnaryClientInterceptor, 0, len(ucc.interceptors)+len(interceptors))
	ins = append(ins, ucc.interceptors...)
	ins = append(ins, interceptors...)

	return UnaryClientChain{interceptors: ins}
}

// Extend extends a chain, by appending another chain to the request flow.
func (ucc UnaryClientChain) Extend(ucchain UnaryClientChain) UnaryClientChain {
	return ucc.Append(ucchain.interceptors...)
}
