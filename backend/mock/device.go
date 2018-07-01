package mock

import (
	"github.com/juzempelde/procwatch/backend"
)

type Device struct {
	DisconnectError error
	IdentifyErrors  map[procwatch.DeviceID]error
}

func (dev *Device) Identify(id procwatch.DeviceID) error {
	if dev.IdentifyErrors == nil {
		return nil
	}
	return dev.Identify(id)
}

func (dev *Device) Disconnect() error {
	return dev.DisconnectError
}
