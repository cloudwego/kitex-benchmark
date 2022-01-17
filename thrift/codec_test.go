// Copyright 2022 ByteDance Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package thrift

import (
	"testing"

	"github.com/cloudwego/frugal"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
)

func makeRequest() *echo.NestedRequest {
	return &echo.NestedRequest{
		NestedStruct: &echo.NestedStruct{
			Struct:     &echo.SimpleStruct{},
			StructList: []*echo.SimpleStruct{{}, {}, {}, {}, {}},
			StructMap:  map[string]*echo.SimpleStruct{"0": {}, "1": {}, "2": {}, "3": {}, "4": {}},
		},
		Request: &echo.Request{},
	}
}

func BenchmarkFrugal(b *testing.B) {
	req := makeRequest()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		size := frugal.EncodedSize(req)
		buf := make([]byte, size)
		_, err := frugal.EncodeObject(buf, nil, req)
		if err != nil {
			b.Fatal(err)
		}
		r := &echo.NestedRequest{}
		_, err = frugal.DecodeObject(buf, r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFastAPI(b *testing.B) {
	req := makeRequest()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		size := req.BLength()
		buf := make([]byte, size)
		req.FastWriteNocopy(buf, nil)
		r := &echo.NestedRequest{}
		_, err := r.FastRead(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFrugalPrealloc(b *testing.B) {
	req := makeRequest()
	size := frugal.EncodedSize(req)
	buf := make([]byte, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := frugal.EncodeObject(buf, nil, req)
		if err != nil {
			b.Fatal(err)
		}
		r := &echo.NestedRequest{}
		_, err = frugal.DecodeObject(buf, r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFastAPIPrealloc(b *testing.B) {
	req := makeRequest()
	size := req.BLength()
	buf := make([]byte, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.FastWriteNocopy(buf, nil)
		r := &echo.NestedRequest{}
		_, err := r.FastRead(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFrugalUnmarshalPrealloc(b *testing.B) {
	req := makeRequest()
	size := frugal.EncodedSize(req)
	buf := make([]byte, size)
	_, err := frugal.EncodeObject(buf, nil, req)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := &echo.NestedRequest{}
		_, err = frugal.DecodeObject(buf, r)
		if err != nil {
			b.Fatal(err)
		}
		if len(r.NestedStruct.StructList) != 5 {
			b.Fatal("wrong data")
		}
	}
}

func BenchmarkFastAPIUnmarshalPrealloc(b *testing.B) {
	req := makeRequest()
	size := req.BLength()
	buf := make([]byte, size)
	req.FastWriteNocopy(buf, nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := &echo.NestedRequest{}
		_, err := r.FastRead(buf)
		if err != nil {
			b.Fatal(err)
		}
		if len(r.NestedStruct.StructList) != 5 {
			b.Fatal("wrong data")
		}
	}
}
