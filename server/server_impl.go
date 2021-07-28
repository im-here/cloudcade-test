package server

import (
	"chat/protocol"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	second = 1
	minute = second * 60
	hours  = minute * 60
	day    = hours * 24
)

type client struct {
	conn       net.Conn
	name       string
	channel    string
	loginAt    int64
	onlineTime int64
	writer     *protocol.CommandWriter
}

type ChatServer struct {
	listener net.Listener
	mutex    *sync.Mutex
	clients  []*client
	message  map[string][]*protocol.Message
}

func NewServer() *ChatServer {
	return &ChatServer{
		mutex:   &sync.Mutex{},
		message: make(map[string][]*protocol.Message),
	}
}

func (s *ChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)

	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", address)

	return err
}

func (s *ChatServer) Close() error {
	return s.listener.Close()
}

func (s *ChatServer) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

func (s *ChatServer) Broadcast(command interface{}, channel string) error {
	if channel != "" {
		for _, client := range s.clients {
			if client.channel == channel {
				client.writer.Write(command)
			}
		}
	} else {
		for _, client := range s.clients {
			client.writer.Write(command)
		}
	}

	return nil
}

func (s *ChatServer) Response(name string, data interface{}) {
	for _, c := range s.clients {
		if c.name == name {
			c.writer.Write(data)
			break
		}
	}
}

func (s *ChatServer) UpdateOnline(name string) {
	for _, c2 := range s.clients {
		if c2.name == name {
			c2.onlineTime = time.Now().Unix() - c2.loginAt
			break
		}
	}
}

func (s *ChatServer) UserStats(name string) string {
	for _, c2 := range s.clients {
		if c2.name == name {
			v := int(c2.onlineTime)
			_day := v / day
			_hour := v % day / hours
			_minute := v % day % hours / minute
			_sec := v % day % hours % minute / second
			return fmt.Sprintf("%dd %dh %dm %ds", _day, _hour, _minute, _sec)
		}
	}
	return "User not online"
}

func (s *ChatServer) SaveMsg(name, msg, channel string) {
	s.message[channel] = append(s.message[channel], &protocol.Message{
		Sender:  name,
		Content: msg,
		Time:    time.Now(),
	})
}

func (s *ChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}

	s.clients = append(s.clients, client)

	return client
}

func (s *ChatServer) remove(client *client) {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *ChatServer) serve(client *client) {
	cmdReader := protocol.NewCommandReader(client.conn)

	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.HeartBeatCommand:
				go s.Response(v.Name, "")
				s.UpdateOnline(v.Name)
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				}, client.channel)
				s.SaveMsg(client.name, v.Message, client.channel)
			case protocol.NameCommand:
				client.name = v.Name
				client.loginAt = time.Now().Unix()
				go s.Response(v.Name, protocol.NameCommand{Name: v.Name})
			case protocol.ChannelCommand:
				client.channel = v.Channel
				go s.Response(v.Name, protocol.ChannelCommand{
					Name:    v.Name,
					Channel: v.Channel,
				})
			case protocol.Command:
				switch v.Cmd {
				case "/popular":
					// TODO
					_ = v.Arg
				case "/stats":
					go s.Response(client.name, protocol.Command{
						Content: s.UserStats(v.Arg),
					})
				}
			case protocol.LoginSuccess:
				// TODO: send chat records
			}
		}

		if err == io.EOF {
			break
		}
	}
}
