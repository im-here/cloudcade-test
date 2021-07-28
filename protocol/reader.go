package protocol

import (
	"bufio"
	"io"
	"log"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read() (interface{}, error) {
	commandName, err := r.reader.ReadString(' ')

	if err != nil {
		return nil, err
	}

	switch commandName {
	case "HEARTBEAT ":
		name, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return HeartBeatCommand{Name: name[:len(name)-1]}, nil
	case "MESSAGE ":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		message, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return MessageCommand{
			Name:    user[:len(user)-1],
			Message: message[:len(message)-1],
		}, nil
	case "SEND ":
		message, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		channel, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return SendCommand{Message: message[:len(message)-1], Channel: channel[:len(channel)-1]}, nil
	case "NAME ":
		name, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return NameCommand{Name: name[:len(name)-1]}, nil
	case "CHANNEL ":
		name, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		channel, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return ChannelCommand{
			Name:    name[:len(name)-1],
			Channel: channel[:len(channel)-1],
		}, nil
	case "CMD ":
		cmd, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		arg, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		content, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return Command{
			Cmd:     cmd[:len(cmd)-1],
			Arg:     arg[:len(arg)-1],
			Content: content[:len(content)-1],
		}, nil
	//case "LOGIN_SUCCESS":
	//	name, err := r.reader.ReadString(' ')
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	msg, err := r.reader.ReadString('\n')
	//	if err != nil {
	//		return nil, err
	//	}
	//	_ = msg
	//	return LoginSuccess{
	//		Name: name[:len(name)-1],
	//		Msg:  nil,
	//	}, nil
	default:
		log.Printf("Unknown command: %v", commandName)
	}

	return nil, UnknownCommand
}

func (r *CommandReader) ReadAll() ([]interface{}, error) {
	var commands []interface{}

	for {
		command, err := r.Read()

		if command != nil {
			commands = append(commands, command)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return commands, err
		}
	}

	return commands, nil
}
