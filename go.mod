module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	dubbo.apache.org/dubbo-go/v3 v3.1.1
	github.com/apache/dubbo-go-hessian2 v1.12.2
	github.com/apache/thrift v0.14.0
	github.com/cloudfoundry/gosigar v1.3.3
	github.com/cloudwego/fastpb v0.0.2
	github.com/cloudwego/kitex v0.4.3
	github.com/dubbogo/grpc-go v1.42.10
	github.com/dubbogo/triple v1.2.2-rc3
	github.com/gogo/protobuf v1.3.2
	github.com/juju/ratelimit v1.0.1
	github.com/lesismal/arpc v1.2.4
	github.com/lesismal/nbio v1.1.23-0.20210815145206-b106d99bce56
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.11
	google.golang.org/grpc v1.52.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/choleraehyq/pid v0.0.18 // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.27.0 // indirect
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
