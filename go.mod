module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	github.com/apache/thrift v0.14.0
	github.com/bytedance/gopkg v0.0.0-20230728082804-614d0af6619b
	github.com/cloudfoundry/gosigar v1.3.3
	github.com/cloudwego/fastpb v0.0.4
	github.com/cloudwego/kitex v0.9.1
	github.com/cloudwego/kitex-tests v0.1.0
	github.com/gogo/protobuf v1.3.2
	github.com/juju/ratelimit v1.0.1
	github.com/lesismal/arpc v1.2.4
	github.com/lesismal/nbio v1.1.23-0.20210815145206-b106d99bce56
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.11
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.28.1
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
