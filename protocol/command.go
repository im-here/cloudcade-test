package protocol

import "errors"

var UnknownCommand = errors.New("UnknownCommand")

type HeartBeatCommand struct {
	Name string
}

type SendCommand struct {
	Message string
	Channel string
}

type NameCommand struct {
	Name string
}

type ChannelCommand struct {
	Name    string
	Channel string
}

type MessageCommand struct {
	Name    string
	Message string
}

type Command struct {
	Cmd     string
	Arg     string
	Content string
}

type LoginSuccess struct {
	Name string
	Msg  []*Message
}
