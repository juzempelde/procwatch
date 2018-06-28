package procwatch_test

import (
	"github.com/juzempelde/procwatch/backend"

	"net"
	"testing"
)

func TestCannotIdentifyTheSameDeviceTwice(t *testing.T) {
	devices := &procwatch.Devices{}

	first := devices.Connect(&net.TCPAddr{
		IP:   net.IP([]byte{10, 0, 0, 24}),
		Port: 13044,
	})
	second := devices.Connect(&net.TCPAddr{
		IP:   net.IP([]byte{10, 0, 0, 88}),
		Port: 13045,
	})

	firstErr := first.Identify(procwatch.DeviceID("megaserver"))
	secondErr := second.Identify(procwatch.DeviceID("megaserver"))

	if firstErr != nil {
		t.Errorf("Expected first identification error to be nil, but got %+v", firstErr)
	}
	if !procwatch.IsAlreadyConnectedError(secondErr) {
		t.Errorf("Expected error to be caused by connecting the same device twice, but got %+v", secondErr)
	}
}

func TestDevicesCanConnectAfterDisconnect(t *testing.T) {
	deviceID := procwatch.DeviceID("foobar")

	devices := &procwatch.Devices{}

	device := devices.Connect(&net.TCPAddr{
		IP:   net.IP([]byte{10, 1, 0, 200}),
		Port: 25000,
	})

	device.Identify(deviceID)
	err := device.Disconnect()

	if err != nil {
		t.Errorf("Expected no error, but got %+v", err)
		return
	}

	device = devices.Connect(&net.TCPAddr{
		IP:   net.IP([]byte{10, 1, 0, 200}),
		Port: 24000,
	})
	err = device.Identify(deviceID)

	if err != nil {
		t.Errorf("Expected no error, but got %+v", err)
	}
}
