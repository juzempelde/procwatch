package rpc

import (
	"github.com/juzempelde/procwatch/backend"
)

type Processes struct{}

func (procs *Processes) Processes(request ProcessesRequest, response *ProcessesResponse) error {
	return nil
}

type ProcessesRequest struct {
	Processes procwatch.Processes
}

type ProcessesResponse struct{}

const processes = "processes"
