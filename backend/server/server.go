package server

import (
	procwatch "github.com/juzempelde/procwatch/backend"
	pwRPC "github.com/juzempelde/procwatch/backend/rpc"

	"fmt"
	"net"
	"time"
)

// Server accepts RPC connections by agents.
type Server struct {
	RPCAddr string

	listener net.Listener

	Devices Devices
}

// Run starts the RPC server.
func (server *Server) Run() error {
	listener, err := net.Listen("tcp", server.RPCAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %+v\n", err)
			// TODO: Better error handling
			continue
		}
		fmt.Printf("Connection from %s\n", conn.RemoteAddr())
		device := server.Devices.Connect(conn.RemoteAddr())
		go func() {
			pwRPC.NewServer(device, pwRPC.RefreshDeadlineByTimeout(time.Now, 5*time.Second, conn.SetDeadline)).ServeConn(conn)
			device.Disconnect()
		}()
	}
}

// Devices abstracts procwatch.Devices.
type Devices interface {
	Connect(addr net.Addr) procwatch.Device
}
