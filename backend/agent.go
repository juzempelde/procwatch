package procwatch

import (
	"context"
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
		var client AgentClient
		client, err = agent.Connector.Connect()
		if err != nil {
			agent.handleError(err)
			time.Sleep(sleepInterval)
			continue
		}
		defer client.Close() // TODO: Handle close error

		for {
			err := client.Identify(context.TODO(), DeviceID(hostID))
			if err != nil {
				agent.handleError(err)
				time.Sleep(sleepInterval)
			} else {
				break
			}
		}

		for {
			processNamesFilter, err := client.ProcessNamesFilter(context.TODO())
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
			agent.handleError(client.Processes(context.TODO(), processes.Filtered(processNamesFilter)))
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
	Connect() (AgentClient, error)
}

type AgentClient interface {
	Close() error
	Identify(ctx context.Context, id DeviceID) error
	ProcessNamesFilter(ctx context.Context) (ProcessFilterNameList, error)
	Processes(ctx context.Context, procs Processes) error
}
