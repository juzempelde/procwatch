package procwatch

import (
	"fmt"
	"sync"
)

type Device struct {
	List *Devices

	ID        DeviceID
	Processes Processes
}

func (dev *Device) Identify(ID string) error {
	return nil
}

type alreadyConnectedError struct {
	ID DeviceID
}

func (err alreadyConnectedError) Error() string {
	return fmt.Sprintf("device %s already connected", err.ID)
}

func (dev *Device) Disconnect() error {
	return nil
}

func (dev *Device) SetProcesses(procs Processes) {}

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

func (dev *Devices) Connect() *Device {
	return nil
}

type DeviceID string

type Processes []*Process

type Process struct {
	Name string
	PID  int
}
