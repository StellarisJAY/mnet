package network

type ServerStartHook func()
type ServerCloseHook func()

type Server interface {
	Start() (error, chan interface{}) // Start server, returns an error if start fail
	Close()                           // Close server

	Protocol() Protocol     // Protocol bind to this connection
	SetProtocol(p Protocol) // Bind a Protocol to this server

	SetServerStartHook(hook ServerStartHook) // Start hook
	SetServerCloseHook(hook ServerCloseHook) // Close hook
}
