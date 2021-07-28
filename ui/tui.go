package tui

import (
	"chat/client"
	"chat/protocol"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"io"
	"strings"
)

var cmdList = []string{
	"/popular",
	"/stats",
}

func StartUi(c client.Client) {
	loginView := NewLoginView()
	channelView := NewChannelView()
	chatView := NewChatView()

	ui, err := tui.New(loginView)
	if err != nil {
		panic(err)
	}

	quit := func() { ui.Quit() }

	ui.SetKeybinding("Esc", quit)
	ui.SetKeybinding("Ctrl+c", quit)

	loginView.OnLogin(func(username string) {
		c.SetName(username)
		ui.SetWidget(channelView)
	})
	channelView.OnChannel(func(channel string) {
		c.SetChannel(channel)
		c.LoginSuccess()
		ui.SetWidget(chatView)
	})

	chatView.OnSubmit(func(msg string) {
		v, cmd := isCmd(msg)
		if v {
			c.SendCommand(cmd)
		} else {
			c.SendMessage(msg)
		}
	})

	go func() {
		for {
			select {
			case err := <-c.Error():

				if err == io.EOF {
					ui.Update(func() {
						chatView.AddMessage("Connection closed connection from server.")
					})
				} else {
					panic(err)
				}
			case msg := <-c.Incoming():
				ui.Update(func() {
					chatView.AddMessage(fmt.Sprintf("%v: %v", msg.Name, msg.Message))
				})
			case msg := <-c.CmdMsg():
				ui.Update(func() {
					chatView.AddMessage(fmt.Sprintf("%v", msg.Content))
				})
			case msg := <-c.Records():
				for _, message := range msg {
					ui.Update(func() {
						chatView.AddMessage(fmt.Sprintf("%v: %v", message.Sender, message.Content))
					})
				}

			}
		}
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

func isCmd(msg string) (bool, protocol.Command) {
	args := strings.Split(msg, " ")
	if len(args) != 2 {
		return false, protocol.Command{}
	}
	cmd := args[0]
	arg := args[1]
	for _, v := range cmdList {
		if cmd == v {
			return true, protocol.Command{
				Cmd: cmd,
				Arg: arg,
			}
		}
	}
	return false, protocol.Command{}
}
