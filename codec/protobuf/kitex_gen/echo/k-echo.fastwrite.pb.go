package echo

import (
	"github.com/cloudwego/kitex/pkg/protocol/bprotoc"
)

func (x *Request) FastWrite(buf []byte) int {
	offset := 0
	if x == nil {
		return offset
	}
	offset += x.fastWriteAction(buf[offset:])
	offset += x.fastWriteMsg(buf[offset:])
	return offset
}

// string
func (x *Request) fastWriteAction(buf []byte) int {
	offset := 0
	if x.Action != "" {
		offset += bprotoc.Binary.WriteString(buf[offset:], 1, x.Action)
	}
	return offset
}

// string
func (x *Request) fastWriteMsg(buf []byte) int {
	offset := 0
	if x.Msg != "" {
		offset += bprotoc.Binary.WriteString(buf[offset:], 2, x.Msg)
	}
	return offset
}
