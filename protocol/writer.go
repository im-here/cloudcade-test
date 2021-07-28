package protocol

import (
	"fmt"
	"io"
)

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

func (w *CommandWriter) writeString(msg string) (err error) {
	_, err = w.writer.Write([]byte(msg))

	return err
}

func (w *CommandWriter) Write(command interface{}) (err error) {
	switch v := command.(type) {
	case HeartBeatCommand:
		err = w.writeString(fmt.Sprintf("HEARTBEAT %v\n", v.Name))
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v %v\n", v.Message, v.Channel))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v\n", v.Name, v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))
	case ChannelCommand:
		err = w.writeString(fmt.Sprintf("CHANNEL %v %v\n", v.Name, v.Channel))
	case Command:
		err = w.writeString(fmt.Sprintf("CMD %v %v %v\n", v.Cmd, v.Arg, v.Content))
	//case LoginSuccess:
	//	err = w.writeString(fmt.Sprintf("LOGIN_SUCCESS %v %v\n", v.Name, v.Msg))
	default:
		err = UnknownCommand
	}

	return err
}
