package server

import (
	procwatch "github.com/juzempelde/procwatch/backend"
	pwRPC "github.com/juzempelde/procwatch/backend/rpc"

	"context"
	"fmt"
	"net"
	"time"
)

// Server accepts RPC connections by agents.
type Server struct {
	RPCAddr string
	Server  *procwatch.Server
}

// Run starts the RPC server.
func (server *Server) Run() error {
	return server.Server.ListenAndServe(&listenAccepter{RPCAddr: server.RPCAddr})
}

type listenAccepter struct {
	RPCAddr  string
	listener net.Listener
}

func (la *listenAccepter) Listen() (procwatch.Accepter, error) {
	var err error
	la.listener, err = net.Listen("tcp", la.RPCAddr)
	if err != nil {
		return nil, err
	}
	return la, nil
}

func (la *listenAccepter) Accept() (procwatch.ServerConnection, error) {
	conn, err := la.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &connection{
		conn: conn,
	}, nil
}

type connection struct {
	conn net.Conn
}

func (conn *connection) Handle(device procwatch.Device) error {
	pwRPC.NewServer(device, pwRPC.RefreshDeadlineByTimeout(time.Now, 5*time.Second, conn.conn.SetDeadline)).ServeConn(conn.conn)
	device.Disconnect()
	return nil // TODO: error handling
}

func (conn *connection) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

func (conn *connection) Shutdown(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
