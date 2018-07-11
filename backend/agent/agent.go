package agent

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/rpc"

	"net"
	"os"
	"time"
)

// Agent runs procwatch as an agent.
type Agent struct {
	ServerRPCAddr string
	ProcessList   ProcessList
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
		processes, err := agent.ProcessList.Current()
		if err != nil {
			continue
		}
		rpcClient.Processes(processes.Filtered(processNamesFilter))
		time.Sleep(sleepInterval)
	}
}

const sleepInterval = time.Second

type ProcessList interface {
	Current() (procwatch.Processes, error)
}

type ProcessListFunc func() (procwatch.Processes, error)

func (f ProcessListFunc) Current() (procwatch.Processes, error) {
	return f()
}
