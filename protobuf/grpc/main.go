/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/cloudwego/kitex-benchmark/perf"
	pb "github.com/cloudwego/kitex-benchmark/protobuf/grpc/grpc_gen"
)

const (
	port = ":8000"
)

var recorder = perf.NewRecorder("GRPC")

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGrpcEchoServer
}

// EchoMsg implements helloworld.GreeterServer
func (s *server) EchoMsg(ctx context.Context, req *pb.GrpcMsg) (resp *pb.GrpcMsg, err error) {
	switch req.Msg {
	case "begin":
		recorder.Begin()
	case "end":
		recorder.End()
		recorder.Report()
	}
	return &pb.GrpcMsg{
		Msg:    req.Msg,
		Finish: req.Finish,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGrpcEchoServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
