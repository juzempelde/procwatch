package procwatch

import (
	"fmt"
	"net"
	"sync"
)

// Device is a device as provided by an agent.
type Device interface {
	Identify(id DeviceID) error
	Disconnect() error
}

type alreadyConnectedError struct {
	ID DeviceID
}

func newAlreadyConnectedError(id DeviceID) *alreadyConnectedError {
	return &alreadyConnectedError{
		ID: id,
	}
}

func (err alreadyConnectedError) Error() string {
	return fmt.Sprintf("device %s already connected", err.ID)
}

// IsAlreadyConnectedError checks wether an error occured because of a device trying to connect twice.
func IsAlreadyConnectedError(err error) bool {
	_, ok := err.(*alreadyConnectedError)
	return ok
}

// Devices contains a bunch of devices, including metadata.
type Devices struct {
	lock sync.Mutex

	currentInternalID int

	byID map[int]*device

	// Possible changes:
	// - Connect device
	// - Identify device
	// - Disconnect device
	// - Set device's processes
	OnChange func()
}

func (devices *Devices) addDevice(dev *device) {
	if devices.byID == nil {
		devices.byID = map[int]*device{}
	}
	devices.byID[dev.internalID] = dev
}

func (devices *Devices) identify(dev *device, id DeviceID) error {
	devices.lock.Lock()
	for _, existingDevice := range devices.byID {
		if existingDevice.id == id && existingDevice.internalID != dev.internalID && existingDevice.connected {
			return newAlreadyConnectedError(id)
		}
	}
	dev.id = id
	devices.lock.Unlock()
	return nil
}

// Connect initially connects a device.
func (devices *Devices) Connect(addr net.Addr) Device {
	devices.lock.Lock()
	devices.currentInternalID++
	newDevice := &device{
		devices:    devices,
		addr:       addr,
		internalID: devices.currentInternalID,
		connected:  true,
	}
	devices.addDevice(newDevice)
	devices.lock.Unlock()
	return newDevice
}

type device struct {
	devices    *Devices
	internalID int

	addr      net.Addr
	id        DeviceID
	connected bool
}

func (dev *device) Identify(id DeviceID) error {
	return dev.devices.identify(dev, id)
}

func (dev *device) Disconnect() error {
	dev.connected = false
	return nil
}

// DeviceID is the self-chosen ID of a device.
type DeviceID string
