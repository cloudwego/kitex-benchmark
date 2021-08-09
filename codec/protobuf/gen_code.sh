# gogo-gen: rpcx, arpc
rm -rf ./gogo_gen && mkdir ./gogo_gen
protoc --gogofaster_out=./gogo_gen ./echo-gogo.proto

# grpc-gen
rm -rf ./grpc_gen && mkdir ./grpc_gen
protoc --go_out=./grpc_gen --go-grpc_out=./grpc_gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ./echo-grpc.proto

# kitex-gen
rm -rf ./kitex_gen && mkdir ./kitex_gen
kitex -type protobuf -module github.com/cloudwego/kitex-benchmark ./echo-kitex.proto
