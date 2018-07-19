package procwatch

import (
	"net"
)

// Observer is informated about changes to devices.
type Observer interface {
	// Init initializes the observer with a device list. After that, all change events are applied against that list.
	Init(list ObservedDeviceList) error

	// Change pushes a difference between the previous and the current state to the observer.
	Change(event ObserverEvent) error
}

// ObservedDeviceList contains observed devices.
type ObservedDeviceList []ObservedDevice

// ObservedDevice represents a single device.
type ObservedDevice struct {
	// DeviceID is the ID of the device. Is nil if the device did not yet identify itself.
	DeviceID   *DeviceID `json:"device_id"`

	// RemoteAddr is the device's remote address. In case of disconnected devices, this is the last known address, so this is
	// never nil.
	RemoteAddr net.Addr  `json:"remote_addr"`

	// Processes is a list of processes the device has reported. In case of disconnected devices, this is the last known list.
	Processes  []Process `json:"processes"`
}

type ObserverEvent interface {
	Type() ObserverEventType
}

type ObserverEventType interface {
	String() string

	hidden()
}
