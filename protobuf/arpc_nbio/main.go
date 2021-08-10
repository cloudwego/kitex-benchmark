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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lesismal/arpc"
	"github.com/lesismal/arpc/codec"
	log "github.com/lesismal/arpc/log"
	"github.com/lesismal/nbio"
	nlog "github.com/lesismal/nbio/logging"

	"github.com/cloudwego/kitex-benchmark/perf"
	pb "github.com/cloudwego/kitex-benchmark/protobuf/arpc/pb_gen"
	"github.com/cloudwego/kitex-benchmark/protobuf/arpc/pbcodec"
)

const (
	port = ":8005"
)

var recorder = perf.NewRecorder("ARPC_NBIO")

func EchoMsg(ctx *arpc.Context) {
	args := &pb.ArpcMsg{}

	if err := ctx.Bind(args); err != nil {
		ctx.Error(err)
		return
	}
	switch args.Msg {
	case "begin":
		recorder.Begin()
	case "end":
		recorder.End()
		recorder.Report()
	}

	reply := &pb.ArpcMsg{
		Msg:    args.Msg,
		Finish: args.Finish,
	}

	ctx.Write(reply)
}

func main() {
	log.SetLevel(log.LevelNone)
	nlog.SetLogger(log.DefaultLogger)

	codec.DefaultCodec = &pbcodec.ProtoBuffer{}

	handler.SetAsyncWrite(false)
	handler.SetAsyncResponse(true)
	handler.Handle("EchoMsg", EchoMsg)

	g := nbio.NewGopher(nbio.Config{
		Network: "tcp",
		Addrs:   []string{port},
	})

	g.OnOpen(onOpen)
	g.OnData(onData)

	err := g.Start()
	if err != nil {
		log.Error("Start failed: %v", err)
	}
	defer g.Stop()

	time.Sleep(time.Second / 10)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

var (
	handler = arpc.NewHandler()
)

// Session .
type Session struct {
	Client *arpc.Client
	Buffer []byte
}

func onOpen(c *nbio.Conn) {
	client := &arpc.Client{Conn: c, Codec: codec.DefaultCodec, Handler: handler}
	session := &Session{
		Client: client,
		Buffer: nil,
	}
	c.SetSession(session)
}

func onData(c *nbio.Conn, data []byte) {
	iSession := c.Session()
	if iSession == nil {
		c.Close()
		return
	}
	session := iSession.(*Session)
	session.Buffer = append(session.Buffer, data...)
	for len(session.Buffer) >= arpc.HeadLen {
		headBuf := session.Buffer[:4]
		header := arpc.Header(headBuf)
		if len(session.Buffer) < arpc.HeadLen+header.BodyLen() {
			return
		}

		msg := &arpc.Message{Buffer: session.Buffer[:arpc.HeadLen+header.BodyLen()]}
		session.Buffer = session.Buffer[arpc.HeadLen+header.BodyLen():]
		handler.OnMessage(session.Client, msg)
	}
}
