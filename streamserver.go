package gonion

import "google.golang.org/grpc"

// NewStreamServer creates a new StreamServer chain.
func NewStreamServer(interceptors ...grpc.StreamServerInterceptor) StreamServerChain {
	return StreamServerChain{append(([]grpc.StreamServerInterceptor)(nil), interceptors...)}
}

// Then chains middlewares and returns the final grpc.StreamServerInterceptor.
func (ssc StreamServerChain) Then(ssi grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	if ssi == nil {
		ssi = func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
			handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}
	}
	newInterceptors := make([]grpc.StreamServerInterceptor, 0, len(ssc.interceptors)+1)
	newInterceptors = append(newInterceptors, ssc.interceptors...)
	newInterceptors = append(newInterceptors, ssi)
	newSSC := StreamServerChain{newInterceptors}

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		return newSSC.chain(srv, stream, info, handler, 0)
	}
}

func (ssc StreamServerChain) chain(srv interface{}, stream grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler, i int) error {
	if i == len(ssc.interceptors) {
		return handler(srv, stream)
	}
	return ssc.interceptors[i](srv, stream, info, func(srv2 interface{}, stream2 grpc.ServerStream) error {
		return ssc.chain(srv2, stream2, info, handler, i+1)
	})
}

// Append extends a chain, by appending interceptors to the request flow.
func (ssc StreamServerChain) Append(interceptors ...grpc.StreamServerInterceptor) StreamServerChain {
	ins := make([]grpc.StreamServerInterceptor, 0, len(ssc.interceptors)+len(interceptors))
	ins = append(ins, ssc.interceptors...)
	ins = append(ins, interceptors...)

	return StreamServerChain{interceptors: ins}
}

// Extend extends a chain, by appending another chain to the request flow.
func (ssc StreamServerChain) Extend(sschain StreamServerChain) StreamServerChain {
	return ssc.Append(sschain.interceptors...)
}
