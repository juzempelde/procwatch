package rpc

import (
	"github.com/juzempelde/procwatch/backend"
)

type Processes struct {
	RefreshDeadline RefreshDeadlineFunc
}

func (procs *Processes) Processes(request ProcessesRequest, response *ProcessesResponse) error {
	refreshDeadline(procs.RefreshDeadline)
	return nil
}

type ProcessesRequest struct {
	Processes procwatch.Processes
}

type ProcessesResponse struct{}

const processes = "processes"
