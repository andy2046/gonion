

# gonion
`import "github.com/andy2046/gonion"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func WithStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption](#WithStreamServer)
* [func WithUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption](#WithUnaryServer)
* [type HTTPChain](#HTTPChain)
  * [func NewHTTP(handlers ...HTTPHandler) HTTPChain](#NewHTTP)
  * [func (hc HTTPChain) Append(handlers ...HTTPHandler) HTTPChain](#HTTPChain.Append)
  * [func (hc HTTPChain) Extend(hchain HTTPChain) HTTPChain](#HTTPChain.Extend)
  * [func (hc HTTPChain) Then(h http.Handler) http.Handler](#HTTPChain.Then)
* [type HTTPHandler](#HTTPHandler)
* [type StreamClientChain](#StreamClientChain)
  * [func NewStreamClient(interceptors ...grpc.StreamClientInterceptor) StreamClientChain](#NewStreamClient)
  * [func (scc StreamClientChain) Append(interceptors ...grpc.StreamClientInterceptor) StreamClientChain](#StreamClientChain.Append)
  * [func (scc StreamClientChain) Extend(scchain StreamClientChain) StreamClientChain](#StreamClientChain.Extend)
  * [func (scc StreamClientChain) Then(ssi grpc.StreamClientInterceptor) grpc.StreamClientInterceptor](#StreamClientChain.Then)
* [type StreamServerChain](#StreamServerChain)
  * [func NewStreamServer(interceptors ...grpc.StreamServerInterceptor) StreamServerChain](#NewStreamServer)
  * [func (ssc StreamServerChain) Append(interceptors ...grpc.StreamServerInterceptor) StreamServerChain](#StreamServerChain.Append)
  * [func (ssc StreamServerChain) Extend(sschain StreamServerChain) StreamServerChain](#StreamServerChain.Extend)
  * [func (ssc StreamServerChain) Then(ssi grpc.StreamServerInterceptor) grpc.StreamServerInterceptor](#StreamServerChain.Then)
* [type UnaryClientChain](#UnaryClientChain)
  * [func NewUnaryClient(interceptors ...grpc.UnaryClientInterceptor) UnaryClientChain](#NewUnaryClient)
  * [func (ucc UnaryClientChain) Append(interceptors ...grpc.UnaryClientInterceptor) UnaryClientChain](#UnaryClientChain.Append)
  * [func (ucc UnaryClientChain) Extend(ucchain UnaryClientChain) UnaryClientChain](#UnaryClientChain.Extend)
  * [func (ucc UnaryClientChain) Then(usi grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor](#UnaryClientChain.Then)
* [type UnaryServerChain](#UnaryServerChain)
  * [func NewUnaryServer(interceptors ...grpc.UnaryServerInterceptor) UnaryServerChain](#NewUnaryServer)
  * [func (usc UnaryServerChain) Append(interceptors ...grpc.UnaryServerInterceptor) UnaryServerChain](#UnaryServerChain.Append)
  * [func (usc UnaryServerChain) Extend(uschain UnaryServerChain) UnaryServerChain](#UnaryServerChain.Extend)
  * [func (usc UnaryServerChain) Then(usi grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor](#UnaryServerChain.Then)
* [type WrappedServerStream](#WrappedServerStream)
  * [func WrapServerStream(stream grpc.ServerStream) WrappedServerStream](#WrapServerStream)
  * [func (w WrappedServerStream) Context() context.Context](#WrappedServerStream.Context)


#### <a name="pkg-files">Package files</a>
[grpc.go](./grpc.go) [http.go](./http.go) [streamclient.go](./streamclient.go) [streamserver.go](./streamserver.go) [unaryclient.go](./unaryclient.go) [unaryserver.go](./unaryserver.go) 





## <a name="WithStreamServer">func</a> [WithStreamServer](./grpc.go?s=1117:1202#L43)
``` go
func WithStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption
```
WithStreamServer returns a grpc.Server config option with multiple stream interceptors.



## <a name="WithUnaryServer">func</a> [WithUnaryServer](./grpc.go?s=864:947#L38)
``` go
func WithUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption
```
WithUnaryServer returns a grpc.Server config option with multiple unary interceptors.




## <a name="HTTPChain">type</a> [HTTPChain](./http.go?s=174:220#L10)
``` go
type HTTPChain struct {
    // contains filtered or unexported fields
}
```
HTTPChain chains http middlewares.







### <a name="NewHTTP">func</a> [NewHTTP](./http.go?s=261:308#L16)
``` go
func NewHTTP(handlers ...HTTPHandler) HTTPChain
```
NewHTTP creates a new HTTP chain.





### <a name="HTTPChain.Append">func</a> (HTTPChain) [Append](./http.go?s=787:848#L37)
``` go
func (hc HTTPChain) Append(handlers ...HTTPHandler) HTTPChain
```
Append extends a chain, by appending handlers to the request flow.




### <a name="HTTPChain.Extend">func</a> (HTTPChain) [Extend](./http.go?s=1081:1135#L46)
``` go
func (hc HTTPChain) Extend(hchain HTTPChain) HTTPChain
```
Extend extends a chain, by appending another chain to the request flow.




### <a name="HTTPChain.Then">func</a> (HTTPChain) [Then](./http.go?s=525:578#L24)
``` go
func (hc HTTPChain) Then(h http.Handler) http.Handler
```
Then chains middlewares and returns the final http.Handler.


	NewHTTP(m1, m2, m3).Then(h)

is equivalent to:


	m1(m2(m3(h)))




## <a name="HTTPHandler">type</a> [HTTPHandler](./http.go?s=89:132#L7)
``` go
type HTTPHandler func(http.Handler) http.Handler
```
HTTPHandler represents a http middleware.










## <a name="StreamClientChain">type</a> [StreamClientChain](./grpc.go?s=134:209#L11)
``` go
type StreamClientChain struct {
    // contains filtered or unexported fields
}
```
StreamClientChain chains grpc StreamClient middlewares.







### <a name="NewStreamClient">func</a> [NewStreamClient](./streamclient.go?s=119:203#L10)
``` go
func NewStreamClient(interceptors ...grpc.StreamClientInterceptor) StreamClientChain
```
NewStreamClient creates a new StreamClient chain.





### <a name="StreamClientChain.Append">func</a> (StreamClientChain) [Append](./streamclient.go?s=1861:1960#L45)
``` go
func (scc StreamClientChain) Append(interceptors ...grpc.StreamClientInterceptor) StreamClientChain
```
Append extends a chain, by appending interceptors to the request flow.




### <a name="StreamClientChain.Extend">func</a> (StreamClientChain) [Extend](./streamclient.go?s=2252:2332#L54)
``` go
func (scc StreamClientChain) Extend(scchain StreamClientChain) StreamClientChain
```
Extend extends a chain, by appending another chain to the request flow.




### <a name="StreamClientChain.Then">func</a> (StreamClientChain) [Then](./streamclient.go?s=378:474#L15)
``` go
func (scc StreamClientChain) Then(ssi grpc.StreamClientInterceptor) grpc.StreamClientInterceptor
```
Then chains middlewares and returns the final grpc.StreamClientInterceptor.




## <a name="StreamServerChain">type</a> [StreamServerChain](./grpc.go?s=272:347#L16)
``` go
type StreamServerChain struct {
    // contains filtered or unexported fields
}
```
StreamServerChain chains grpc StreamServer middlewares.







### <a name="NewStreamServer">func</a> [NewStreamServer](./streamserver.go?s=102:186#L6)
``` go
func NewStreamServer(interceptors ...grpc.StreamServerInterceptor) StreamServerChain
```
NewStreamServer creates a new StreamServer chain.





### <a name="StreamServerChain.Append">func</a> (StreamServerChain) [Append](./streamserver.go?s=1514:1613#L40)
``` go
func (ssc StreamServerChain) Append(interceptors ...grpc.StreamServerInterceptor) StreamServerChain
```
Append extends a chain, by appending interceptors to the request flow.




### <a name="StreamServerChain.Extend">func</a> (StreamServerChain) [Extend](./streamserver.go?s=1905:1985#L49)
``` go
func (ssc StreamServerChain) Extend(sschain StreamServerChain) StreamServerChain
```
Extend extends a chain, by appending another chain to the request flow.




### <a name="StreamServerChain.Then">func</a> (StreamServerChain) [Then](./streamserver.go?s=361:457#L11)
``` go
func (ssc StreamServerChain) Then(ssi grpc.StreamServerInterceptor) grpc.StreamServerInterceptor
```
Then chains middlewares and returns the final grpc.StreamServerInterceptor.




## <a name="UnaryClientChain">type</a> [UnaryClientChain](./grpc.go?s=408:481#L21)
``` go
type UnaryClientChain struct {
    // contains filtered or unexported fields
}
```
UnaryClientChain chains grpc UnaryClient middlewares.







### <a name="NewUnaryClient">func</a> [NewUnaryClient](./unaryclient.go?s=117:198#L10)
``` go
func NewUnaryClient(interceptors ...grpc.UnaryClientInterceptor) UnaryClientChain
```
NewUnaryClient creates a new UnaryClient chain.





### <a name="UnaryClientChain.Append">func</a> (UnaryClientChain) [Append](./unaryclient.go?s=1804:1900#L45)
``` go
func (ucc UnaryClientChain) Append(interceptors ...grpc.UnaryClientInterceptor) UnaryClientChain
```
Append extends a chain, by appending interceptors to the request flow.




### <a name="UnaryClientChain.Extend">func</a> (UnaryClientChain) [Extend](./unaryclient.go?s=2190:2267#L54)
``` go
func (ucc UnaryClientChain) Extend(ucchain UnaryClientChain) UnaryClientChain
```
Extend extends a chain, by appending another chain to the request flow.




### <a name="UnaryClientChain.Then">func</a> (UnaryClientChain) [Then](./unaryclient.go?s=370:463#L15)
``` go
func (ucc UnaryClientChain) Then(usi grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor
```
Then chains middlewares and returns the final grpc.UnaryClientInterceptor.




## <a name="UnaryServerChain">type</a> [UnaryServerChain](./grpc.go?s=542:615#L26)
``` go
type UnaryServerChain struct {
    // contains filtered or unexported fields
}
```
UnaryServerChain chains grpc UnaryServer middlewares.







### <a name="NewUnaryServer">func</a> [NewUnaryServer](./unaryserver.go?s=117:198#L10)
``` go
func NewUnaryServer(interceptors ...grpc.UnaryServerInterceptor) UnaryServerChain
```
NewUnaryServer creates a new UnaryServer chain.





### <a name="UnaryServerChain.Append">func</a> (UnaryServerChain) [Append](./unaryserver.go?s=1553:1649#L44)
``` go
func (usc UnaryServerChain) Append(interceptors ...grpc.UnaryServerInterceptor) UnaryServerChain
```
Append extends a chain, by appending interceptors to the request flow.




### <a name="UnaryServerChain.Extend">func</a> (UnaryServerChain) [Extend](./unaryserver.go?s=1939:2016#L53)
``` go
func (usc UnaryServerChain) Extend(uschain UnaryServerChain) UnaryServerChain
```
Extend extends a chain, by appending another chain to the request flow.




### <a name="UnaryServerChain.Then">func</a> (UnaryServerChain) [Then](./unaryserver.go?s=370:463#L15)
``` go
func (usc UnaryServerChain) Then(usi grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor
```
Then chains middlewares and returns the final grpc.UnaryServerInterceptor.




## <a name="WrappedServerStream">type</a> [WrappedServerStream](./grpc.go?s=687:771#L31)
``` go
type WrappedServerStream struct {
    grpc.ServerStream
    WrappedContext context.Context
}
```
WrappedServerStream is a wrapped grpc.ServerStream with context.







### <a name="WrapServerStream">func</a> [WrapServerStream](./grpc.go?s=1499:1566#L53)
``` go
func WrapServerStream(stream grpc.ServerStream) WrappedServerStream
```
WrapServerStream returns a ServerStream with the ability to overwrite context.





### <a name="WrappedServerStream.Context">func</a> (WrappedServerStream) [Context](./grpc.go?s=1332:1386#L48)
``` go
func (w WrappedServerStream) Context() context.Context
```
Context returns the wrapper's WrappedContext.







