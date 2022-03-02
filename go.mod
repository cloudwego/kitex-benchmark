module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	github.com/apache/thrift v0.14.0
	github.com/cloudwego/kitex v0.1.4-0.20220120104459-2389480e8498
	github.com/gogo/protobuf v1.3.2
	github.com/lesismal/arpc v1.2.4
	github.com/lesismal/nbio v1.1.23-0.20210815145206-b106d99bce56
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.11
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

replace github.com/cloudwego/kitex => github.com/sinnera/kitex v0.0.4-0.20220216092555-831e8f34cd33

replace google.golang.org/grpc => github.com/sinnera/grpc-go v1.45.0-dev.0.20220216092540-33f7497c01d8
