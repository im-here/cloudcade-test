package client

import (
	"log"
	"net"
	"time"

	"chat/protocol"
)

type ChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	channel   string
	error     chan error
	incoming  chan protocol.MessageCommand
	cmdMsg    chan protocol.Command
	records   chan []*protocol.Message
}

func NewClient() *ChatClient {
	return &ChatClient{
		incoming: make(chan protocol.MessageCommand),
		cmdMsg:   make(chan protocol.Command),
		records:  make(chan []*protocol.Message),
		error:    make(chan error),
	}
}

func (c *ChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
		c.cmdReader = protocol.NewCommandReader(conn)
		c.cmdWriter = protocol.NewCommandWriter(conn)
	}

	return err
}

func (c *ChatClient) Start() {
	c.HeartBeat()
	for {
		cmd, err := c.cmdReader.Read()

		if err != nil {
			c.error <- err
			break
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			case protocol.Command:
				c.cmdMsg <- v
			case protocol.NameCommand:
				c.name = v.Name
			case protocol.ChannelCommand:
				c.channel = v.Channel
			case protocol.LoginSuccess:
				c.records <- v.Msg
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

func (c *ChatClient) HeartBeat() error {
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-ticker.C:
				c.Send(protocol.HeartBeatCommand{Name: c.name})
			}

		}
	}()
	return nil
}

func (c *ChatClient) Close() error {
	return c.conn.Close()
}

func (c *ChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *ChatClient) CmdMsg() chan protocol.Command {
	return c.cmdMsg
}

func (c *ChatClient) Records() chan []*protocol.Message {
	return c.records
}

func (c *ChatClient) Error() chan error {
	return c.error
}

func (c *ChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *ChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{Name: name})
}

func (c *ChatClient) GetName() string {
	return c.name
}

func (c *ChatClient) LoginSuccess() error {
	return c.Send(protocol.LoginSuccess{Name: c.name})
}

func (c *ChatClient) SetChannel(channel string) error {
	return c.Send(protocol.ChannelCommand{Name: c.name, Channel: channel})
}

func (c *ChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
		Channel: c.channel,
	})
}

func (c *ChatClient) SendCommand(command protocol.Command) error {
	return c.Send(command)
}
