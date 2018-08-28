package procwatch

import (
	"context"
	"fmt"
	"net"
)

// Server is an abstract representation of a running server.
type Server struct {
	Devices *Devices
}

// ListenAndServe opens a listener and starts accepting connections from agents.
func (srv *Server) ListenAndServe(listener Listener) error {
	accepter, err := listener.Listen()
	if err != nil {
		return err
	}
	for {
		conn, err := accepter.Accept()
		if err != nil {
			fmt.Printf("Accept error: %+v\n", err)
			// TODO: Better error handling
			continue
		}
		fmt.Printf("Connection from %s\n", conn.RemoteAddr())
		device := srv.Devices.Connect(conn.RemoteAddr())
		go func() {
			conn.Handle(device)
			device.Disconnect()
		}()
	}
}

type Listener interface {
	// Listen starts listening.
	Listen() (Accepter, error)
}

type Accepter interface {
	// Accept returns a connection initiated from the outside.
	Accept() (ServerConnection, error)
}

// ServerConnection represents a connection from an agent.
type ServerConnection interface {
	// RemoteAddr returns the agent's remote address.
	RemoteAddr() net.Addr

	// Shutdown closes the connection.
	Shutdown(context.Context) error

	// Handle lets the connection handle commands from the agent.
	Handle(Device) error
}
