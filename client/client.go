package client

import "chat/protocol"

//go:generate mockgen -destination=mock/client.go -package=servermock client Client
type Client interface {
	Dial(address string) error
	Start()
	Close() error
	HeartBeat() error
	Send(command interface{}) error
	SetName(name string) error
	GetName() string
	SetChannel(channel string) error
	LoginSuccess() error
	SendMessage(message string) error
	SendCommand(command protocol.Command) error
	Error() chan error
	Incoming() chan protocol.MessageCommand
	CmdMsg() chan protocol.Command
	Records() chan []*protocol.Message
}
