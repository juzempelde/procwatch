package server

import (
	"net"
	"net/rpc"
)

// Server accepts RPC connections by agents.
type Server struct {
	RPCAddr string
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
