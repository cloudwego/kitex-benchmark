package pkg

import "context"

type Message struct {
	Action string `hessian:"Action"`
	Msg    string `hessian:"Msg"`
}

func (u Message) String() string {
	return u.String()
}

func (u *Message) JavaClassName() string {
	return "org.Request"
}

type EchoProvider struct {
	Echo func(ctx context.Context, request *Message) (*Message, error)
}
