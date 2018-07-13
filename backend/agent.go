package procwatch

import (
	"time"
)

type Agent struct {
	Connector      Connector
	ProcessList    ProcessList
	HostIDProvider HostIDProvider
	ErrorHandler   func(error)
}

func (agent *Agent) Run() error {
	hostID := ""
	var err error
	for {
		hostID, err = agent.HostIDProvider.HostID()
		if err != nil {
			agent.handleError(err)
			time.Sleep(sleepInterval)
		}
	}

	for {
		var client Client
		client, err = agent.Connector.Connect()
		if err != nil {
			agent.handleError(err)
			time.Sleep(sleepInterval)
			continue
		}
		defer client.Close() // TODO: Handle close error

		for {
			err := client.Identify(DeviceID(hostID))
			if err != nil {
				agent.handleError(err)
				time.Sleep(sleepInterval)
			}
		}

		for {
			processNamesFilter, err := client.ProcessNamesFilter()
			if err != nil {
				agent.handleError(err)
				time.Sleep(sleepInterval)
				continue
			}

			processes, err := agent.ProcessList.Current()
			if err != nil {
				agent.handleError(err)
				time.Sleep(sleepInterval)
				continue
			}
			agent.handleError(client.Processes(processes.Filtered(processNamesFilter)))
			time.Sleep(sleepInterval)
		}
	}
}

func (agent *Agent) handleError(err error) {
	if handle := agent.ErrorHandler; handle != nil {
		handle(err)
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
