package gonion

import "net/http"

type (
	// HTTPHandler represents a http middleware.
	HTTPHandler func(http.Handler) http.Handler

	// HTTPChain chains http middlewares.
	HTTPChain struct {
		handlers []HTTPHandler
	}
)

// NewHTTP creates a new HTTP chain.
func NewHTTP(handlers ...HTTPHandler) HTTPChain {
	return HTTPChain{handlers: append(([]HTTPHandler)(nil), handlers...)}
}

// Then chains middlewares and returns the final http.Handler.
//     NewHTTP(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
func (hc HTTPChain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range hc.handlers {
		h = hc.handlers[len(hc.handlers)-1-i](h)
	}

	return h
}

// Append extends a chain, by appending handlers to the request flow.
func (hc HTTPChain) Append(handlers ...HTTPHandler) HTTPChain {
	h := make([]HTTPHandler, 0, len(hc.handlers)+len(handlers))
	h = append(h, hc.handlers...)
	h = append(h, handlers...)

	return HTTPChain{handlers: h}
}

// Extend extends a chain, by appending another chain to the request flow.
func (hc HTTPChain) Extend(hchain HTTPChain) HTTPChain {
	return hc.Append(hchain.handlers...)
}
