module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	github.com/apache/thrift v0.14.0
	github.com/cloudfoundry/gosigar v1.3.3
	github.com/cloudwego/fastpb v0.0.4
	github.com/cloudwego/kitex v0.9.3-rc
	github.com/gogo/protobuf v1.3.2
	github.com/lesismal/arpc v1.2.4
	github.com/lesismal/nbio v1.1.23-0.20210815145206-b106d99bce56
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.11
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.28.1
)

replace (
	github.com/cloudwego/kitex => github.com/cloudwego/kitex v0.9.3-rc2.0.20240523061308-581cc2e538d4
	github.com/cloudwego/netpoll => github.com/cloudwego/netpoll v0.6.1-0.20240516030022-a9a224c3e494
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
