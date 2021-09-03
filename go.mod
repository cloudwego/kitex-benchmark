module github.com/cloudwego/kitex-benchmark

go 1.15

require (
	github.com/apache/thrift v0.14.0
	github.com/cloudwego/kitex v0.0.5-0.20210902124705-e0bb1133da2b
	github.com/gogo/protobuf v1.3.2
	github.com/lesismal/arpc v1.1.9
	github.com/lesismal/nbio v1.1.23-0.20210815145206-b106d99bce56
	github.com/montanaflynn/stats v0.6.6
	github.com/smallnest/rpcx v1.6.7
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210811021853-ddbe55d93216 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/grpc/examples v0.0.0-20210722024238-c513103bee39 // indirect
	google.golang.org/protobuf v1.27.1
)

replace github.com/apache/thrift => github.com/apache/thrift v0.14.2

replace github.com/cloudwego/kitex => github.com/sinnera/kitex v0.0.4-0.20210909101723-fd3cc8b8ae62
