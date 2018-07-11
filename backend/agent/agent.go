package agent

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/rpc"

	"net"
	"os"
	"time"

	"github.com/mitchellh/go-ps"
)

// Agent runs procwatch as an agent.
type Agent struct {
	ServerRPCAddr string
}

// Run connects to the server's address via RPC, registers its ID and sends process informations.
func (agent *Agent) Run() error {
	hostID, err := os.Hostname()
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", agent.ServerRPCAddr)
	if err != nil {
		return err
	}
	rpcClient := rpc.NewClient(conn)
	defer rpcClient.Close() // TODO: Handle close error

	err = rpcClient.Identify(procwatch.DeviceID(hostID))
	if err != nil {
		return err
	}

	processNamesFilter, err := rpcClient.ProcessNamesFilter()
	if err != nil {
		return err
	}

	for {
		processes, err := processes()
		if err != nil {
			continue
		}
		processes = processes.Filtered(processNamesFilter)
		time.Sleep(sleepInterval)
	}
}

const sleepInterval = time.Second

func processes() (procwatch.Processes, error) {
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
