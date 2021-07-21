module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	github.com/apache/thrift v0.14.0
	github.com/cloudwego/kitex v0.0.2-0.20210729075918-8053707e7e0a
	github.com/gogo/protobuf v1.3.2
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.4
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/examples v0.0.0-20210722024238-c513103bee39 // indirect
	google.golang.org/protobuf v1.27.1
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
