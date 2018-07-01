package procwatch

import (
	"fmt"
	"net"
	"sync"
)

// Device is a device as provided by an agent.
type Device interface {
	// Identify identifies the device with an error. Returns an error if the device already is connected.
	Identify(id DeviceID) error

	// Disconnects the device.
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

func (devices *Devices) withLock(f func()) {
	devices.lock.Lock()
	f()
	devices.lock.Unlock()
}

func (devices *Devices) addDevice(dev *device) {
	if devices.byID == nil {
		devices.byID = map[int]*device{}
	}
	devices.byID[dev.internalID] = dev
}

// Connect initially connects a device.
func (devices *Devices) Connect(addr net.Addr) Device {
	devices.lock.Lock()
	devices.currentInternalID++
	newDevice := &device{
		devices:    devices,
		addr:       addr,
		internalID: devices.currentInternalID,
		status:     pending,
	}
	devices.addDevice(newDevice)
	devices.lock.Unlock()
	return newDevice
}

type device struct {
	devices    *Devices
	internalID int
	status     deviceStatus

	addr net.Addr
	id   DeviceID
}

type deviceStatus interface {
	identify(*Devices, *device, DeviceID) error
	disconnect(*device) error
}

func (dev *device) Identify(id DeviceID) error {
	var err error
	dev.devices.withLock(
		func() {
			err = dev.status.identify(dev.devices, dev, id)
		},
	)
	return err
}

func (dev *device) Disconnect() error {
	var err error
	dev.devices.withLock(
		func() {
			err = dev.status.disconnect(dev)
		},
	)
	return err
}

// DeviceID is the self-chosen ID of a device.
type DeviceID string

// statusPending implements deviceStatus and represents a device which is connected, but not yet identified.
type statusPending struct{}

func (status *statusPending) identify(devs *Devices, dev *device, id DeviceID) error {
	for _, existingDevice := range devs.byID {
		if existingDevice.id == id && existingDevice.internalID != dev.internalID && existingDevice.status != disconnected {
			return newAlreadyConnectedError(id)
		}
	}
	dev.status = connected
	dev.id = id
	return nil
}

func (status *statusPending) disconnect(dev *device) error {
	dev.status = disconnected
	return nil
}

var pending = &statusPending{}

// statusConnected implements deviceStatus and represents a device which is connected and identified.
type statusConnected struct{}

func (status *statusConnected) identify(*Devices, *device, DeviceID) error {
	return fmt.Errorf("already identified")
}

func (status *statusConnected) disconnect(dev *device) error {
	dev.status = disconnected
	return nil
}

var connected = &statusConnected{}

// statusDisconnected implements deviceStatus and represents a device which has disconnected.
type statusDisconnected struct{}

func (status *statusDisconnected) identify(*Devices, *device, DeviceID) error {
	return fmt.Errorf("cannot identify after disconnect")
}

func (status *statusDisconnected) disconnect(*device) error {
	return fmt.Errorf("already disconnected")
}

var disconnected = &statusDisconnected{}
