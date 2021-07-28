package main

import "chat/server"

func main() {
	var s *server.ChatServer
	s = server.NewServer()
	err := s.Listen(":8080")
	if err != nil {
		panic(err)
	}

	// start the server
	s.Start()
}
