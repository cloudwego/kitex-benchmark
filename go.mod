module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	github.com/apache/thrift v0.14.0
	github.com/cloudwego/kitex v0.0.5-0.20210924035959-bb859c633696
	github.com/gogo/protobuf v1.3.2
	github.com/lesismal/arpc v1.1.9
	github.com/lesismal/nbio v1.1.23-0.20210815145206-b106d99bce56
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.7
	google.golang.org/grpc v1.39.1
	google.golang.org/grpc/examples v0.0.0-20210923214018-6ff68b489ecb // indirect
	google.golang.org/protobuf v1.27.1
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
