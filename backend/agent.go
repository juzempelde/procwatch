package procwatch

import (
	"time"
)

type Agent struct {
	Connector      Connector
	ProcessList    ProcessList
	HostIDProvider HostIDProvider
}

func (agent *Agent) Run() error {
	hostID, err := agent.HostIDProvider.HostID()
	if err != nil {
		return err
	}

	rpcClient, err := agent.Connector.Connect()
	if err != nil {
		return err
	}
	defer rpcClient.Close() // TODO: Handle close error

	err = rpcClient.Identify(DeviceID(hostID))
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
	Current() (Processes, error)
}

type ProcessListFunc func() (Processes, error)

func (f ProcessListFunc) Current() (Processes, error) {
	return f()
}

type HostIDProvider interface {
	HostID() (string, error)
}

type HostIDProviderFunc func() (string, error)

func (f HostIDProviderFunc) HostID() (string, error) {
	return f()
}

type Connector interface {
	Connect() (Client, error)
}

type Client interface {
	Close() error
	Identify(id DeviceID) error
	ProcessNamesFilter() (ProcessFilterNameList, error)
	Processes(procs Processes) error
}
