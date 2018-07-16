package rpc

import (
	"github.com/juzempelde/procwatch/backend"

	"net/rpc"
	"time"
)

// NewServer creates a new RPC server for a single device.
func NewServer(device procwatch.Device, refreshDeadline RefreshDeadlineFunc) *rpc.Server {
	srv := rpc.NewServer()
	srv.RegisterName(identificator, &DeviceIdentification{Device: device, RefreshDeadline: refreshDeadline})
	srv.RegisterName(processFilter, &ProcessNameFilter{Names: []string{}, RefreshDeadline: refreshDeadline})
	srv.RegisterName(processes, &Processes{RefreshDeadline: refreshDeadline})
	return srv
}

type RefreshDeadlineFunc func() error

func refreshDeadline(f RefreshDeadlineFunc) error {
	if f != nil {
		return f()
	}
	return nil
}

type SetDeadlineFunc func(time.Time) error

func RefreshDeadlineByTimeout(now func() time.Time, duration time.Duration, setDeadline SetDeadlineFunc) RefreshDeadlineFunc {
	return func() error {
		return setDeadline(now().Add(duration))
	}
}
