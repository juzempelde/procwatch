package server

import (
	procwatch "github.com/juzempelde/procwatch/backend"

	"net"
	"net/rpc"
)

// Server accepts RPC connections by agents.
type Server struct {
	RPCAddr string

	listener net.Listener

	Devices Devices
}

// Run starts the RPC server.
func (server *Server) Run() error {
	rpcServer := server.createRPCServer()
	listener, err := net.Listen("tcp", server.RPCAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			// What to do here?
			continue
		}
		go rpcServer.ServeConn(conn) // TODO: Shutdown
	}
}

func (server *Server) createRPCServer() *rpc.Server {
	rpcServer := rpc.NewServer()
	return rpcServer
}

type Devices interface {
	Connect(addr net.Addr) procwatch.Device
}
