package procwatch

import (
	"fmt"
	"net"
	"sync"
)

// Device is a device as provided by an agent.
type Device interface{}

type alreadyConnectedError struct {
	ID DeviceID
}

func (err alreadyConnectedError) Error() string {
	return fmt.Sprintf("device %s already connected", err.ID)
}

type Devices struct {
	PendingDevices []*Devices
	DevicesByID    map[DeviceID]*Device

	Lock sync.Mutex

	// Possible changes:
	// - Connect device
	// - Identify device
	// - Disconnect device
	// - Set device's processes
	OnChange func()
}

// Connect initially connects a device.
func (devices *Devices) Connect(addr net.Addr) Device {
	return &device{
		devices: devices,
		addr:    addr,
	}
}

type device struct {
	devices *Devices

	addr net.Addr
}

type DeviceID string
