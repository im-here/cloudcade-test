package protocol

import "time"

type Message struct {
	Sender  string
	Content string
	Time    time.Time
}
