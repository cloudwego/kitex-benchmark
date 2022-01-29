package echo

import (
	"github.com/cloudwego/kitex/pkg/protocol/bprotoc"
)

func (x *Request) Size() int {
	l := 0
	if x == nil {
		return l
	}
	l += x.sizeAction()
	l += x.sizeMsg()
	return l
}

func (x *Request) sizeAction() (n int) {
	if x.Action != "" {
		n += bprotoc.Binary.SizeString(14, x.Action)
	}
	return n
}

func (x *Request) sizeMsg() (n int) {
	if x.Msg != "" {
		n += bprotoc.Binary.SizeString(14, x.Msg)
	}
	return n
}

func (x *Response) Size() int {
	l := 0
	if x == nil {
		return l
	}
	l += x.sizeAction()
	l += x.sizeMsg()
	return l
}

func (x *Response) sizeAction() (n int) {
	if x.Action != "" {
		n += bprotoc.Binary.SizeString(14, x.Action)
	}
	return n
}

func (x *Response) sizeMsg() (n int) {
	if x.Msg != "" {
		n += bprotoc.Binary.SizeString(14, x.Msg)
	}
	return n
}
