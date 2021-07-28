package server

//go:generate mockgen -destination=mock/server.go -package=servermock server Server
type Server interface {
	Listen(address string) error
	Broadcast(command interface{}, channel string) error
	Response(name string, data interface{})
	UpdateOnline(name string)
	UserStats(name string) string
	SaveMsg(name, msg, channel string)
	Start()
	Close() error
}
