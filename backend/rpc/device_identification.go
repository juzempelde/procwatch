package rpc

import (
	"github.com/juzempelde/procwatch/backend"
)

type DeviceIdentification struct {
	Device          procwatch.Device
	RefreshDeadline RefreshDeadlineFunc
}

func (rpc *DeviceIdentification) Identify(request *IdentificationRequest, response *IdentificationResponse) error {
	refreshDeadline(rpc.RefreshDeadline)
	err := rpc.Device.Identify(request.ID)
	if err != nil {
		response.ErrReason = err.Error()
	}
	return nil
}

type IdentificationRequest struct {
	ID procwatch.DeviceID
}

type IdentificationResponse struct {
	ErrReason string
}

const identificator = "identificator"
