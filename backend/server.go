package procwatch

import (
	"context"
	"fmt"
	"net"
)

type Server struct {
	Devices *Devices
}

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
	Listen() (Accepter, error)
}

type Accepter interface {
	Accept() (ServerConnection, error)
}

type ServerConnection interface {
	RemoteAddr() net.Addr
	Shutdown(context.Context) error
	Handle(Device) error
}
