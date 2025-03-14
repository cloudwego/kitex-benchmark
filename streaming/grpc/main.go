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
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpcg "github.com/cloudwego/kitex-benchmark/codec/protobuf/grpc_gen"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8000
)

var recorder = perf.NewRecorder("GRPC@Server")

type server struct {
	grpcg.UnimplementedSEchoServer
}

func (s *server) Echo(stream grpcg.SEcho_EchoServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	if md == nil || len(md["header"]) == 0 || md["header"][0] != "hello" {
		return fmt.Errorf("invalid header: %v", md)
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		action, msg := runner.ProcessRequest(recorder, req.Action, req.Msg)
		err = stream.Send(&grpcg.Response{
			Action: action,
			Msg:    msg,
		})
		if err != nil {
			log.Printf("stream send failed: %v\n", err)
			return err
		}
	}
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcg.RegisterSEchoServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
