package tui

import "github.com/marcusolsson/tui-go"

type ChannelHandler func(string)

type ChannelView struct {
	tui.Box
	frame          *tui.Box
	ChannelHandler ChannelHandler
}

func NewChannelView() *ChannelView {
	user := tui.NewEntry()
	user.SetFocused(true)
	user.SetSizePolicy(tui.Maximum, tui.Maximum)

	label := tui.NewLabel("Enter your channel(Press enter to enter the public channel): ")
	user.SetSizePolicy(tui.Expanding, tui.Maximum)

	userBox := tui.NewHBox(
		label,
		user,
	)
	userBox.SetBorder(true)
	userBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	view := &ChannelView{}
	view.frame = tui.NewVBox(
		tui.NewSpacer(),
		tui.NewPadder(-4, 0, tui.NewPadder(4, 0, userBox)),
		tui.NewSpacer(),
	)
	view.Append(view.frame)

	user.OnSubmit(func(e *tui.Entry) {
		if view.ChannelHandler != nil {
			view.ChannelHandler(e.Text())
		}
		e.SetText("")
	})

	return view
}

func (v *ChannelView) OnChannel(handler ChannelHandler) {
	v.ChannelHandler = handler
}
