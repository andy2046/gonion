package gonion

import (
	"context"

	"google.golang.org/grpc"
)

// NewUnaryServer creates a new UnaryServer chain.
func NewUnaryServer(interceptors ...grpc.UnaryServerInterceptor) UnaryServerChain {
	return UnaryServerChain{append(([]grpc.UnaryServerInterceptor)(nil), interceptors...)}
}

// Then chains middlewares and returns the final grpc.UnaryServerInterceptor.
func (usc UnaryServerChain) Then(usi grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	if usi == nil {
		usi = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}
	}
	newInterceptors := make([]grpc.UnaryServerInterceptor, 0, len(usc.interceptors)+1)
	newInterceptors = append(newInterceptors, usc.interceptors...)
	newInterceptors = append(newInterceptors, usi)
	newUSC := UnaryServerChain{newInterceptors}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		return newUSC.chain(ctx, req, info, handler, 0)
	}
}

func (usc UnaryServerChain) chain(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler, i int) (interface{}, error) {
	if i == len(usc.interceptors) {
		return handler(ctx, req)
	}
	return usc.interceptors[i](ctx, req, info, func(ctx2 context.Context, req2 interface{}) (interface{}, error) {
		return usc.chain(ctx2, req2, info, handler, i+1)
	})
}

// Append extends a chain, by appending interceptors to the request flow.
func (usc UnaryServerChain) Append(interceptors ...grpc.UnaryServerInterceptor) UnaryServerChain {
	ins := make([]grpc.UnaryServerInterceptor, 0, len(usc.interceptors)+len(interceptors))
	ins = append(ins, usc.interceptors...)
	ins = append(ins, interceptors...)

	return UnaryServerChain{interceptors: ins}
}

// Extend extends a chain, by appending another chain to the request flow.
func (usc UnaryServerChain) Extend(uschain UnaryServerChain) UnaryServerChain {
	return usc.Append(uschain.interceptors...)
}
