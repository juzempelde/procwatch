package server

import (
	procwatch "github.com/juzempelde/procwatch/backend"

	"net"
	"net/rpc"
)

// Server accepts RPC connections by agents.
type Server struct {
	RPCAddr string

	Devices Devices
}

// Run starts the RPC server.
func (server *Server) Run() error {
	listener, err := net.Listen("tcp", server.RPCAddr)
	if err != nil {
		return err
	}
	rpcServer := rpc.NewServer()
	rpcServer.Accept(listener)
	return nil
}

type Devices interface {
	Connect(addr net.Addr) procwatch.Device
}
