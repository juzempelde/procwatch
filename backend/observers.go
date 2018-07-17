package procwatch

import (
	"net"
)

type Observer interface {
	Init(list ObservedDeviceList) error
	Change(event ObserverEvent) error
}

type ObservedDeviceList []ObservedDevice

type ObservedDevice struct {
	DeviceID   *DeviceID `json:"device_id"`
	RemoteAddr net.Addr  `json:"remote_addr"`
	Processes  []Process `json:"processes"`
}

type ObserverEvent interface {
	Type() ObserverEventType
}

type ObserverEventType interface {
	String() string

	hidden()
}
