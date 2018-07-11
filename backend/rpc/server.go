package rpc

import (
	"github.com/juzempelde/procwatch/backend"

	"net/rpc"
)

// NewServer creates a new RPC server for a single device.
func NewServer(device procwatch.Device) *rpc.Server {
	srv := rpc.NewServer()
	srv.RegisterName(identificator, &DeviceIdentification{Device: device})
	srv.RegisterName(processFilter, &ProcessNameFilter{Names: make([]string, 0)})
	return srv
}
