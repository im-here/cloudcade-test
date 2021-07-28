package main

import (
	"chat/client"
	tui "chat/ui"
	"flag"
	"log"
)

func main() {
	address := flag.String("server", "localhost:8080", "Which server to connect to")

	flag.Parse()

	c := client.NewClient()
	err := c.Dial(*address)

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	// start client
	go c.Start()

	tui.StartUi(c)
}
