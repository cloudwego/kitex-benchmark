package echo

import (
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/pkg/protocol/bprotoc"
)

func (x *Request) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadAction(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadMsg(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = bprotoc.Binary.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, errors.New("cannot parse invalid wire-format data")
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_Request[number], err)
}

func (x *Request) fastReadAction(buf []byte, _type int8) (int, error) {
	offset := 0
	v, l, err := bprotoc.Binary.ReadString(buf[offset:], _type)
	if err != nil {
		return offset, err
	}
	offset += l
	x.Action = v
	return offset, nil
}

func (x *Request) fastReadMsg(buf []byte, _type int8) (int, error) {
	offset := 0
	v, l, err := bprotoc.Binary.ReadString(buf[offset:], _type)
	if err != nil {
		return offset, err
	}
	offset += l
	x.Msg = v
	return offset, nil
}
