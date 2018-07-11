package ps

import (
	"github.com/juzempelde/procwatch/backend"

	"github.com/mitchellh/go-ps"
)

// CurrentProcesses returns a list of currently running processes.
func CurrentProcesses() (procwatch.Processes, error) {
	psProcs, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	procs := make(procwatch.Processes, len(psProcs))
	for _, psProc := range psProcs {
		procs = append(
			procs,
			&procwatch.Process{
				PID:  psProc.Pid(),
				Name: psProc.Executable(),
			},
		)
	}
	return procs, nil
}
