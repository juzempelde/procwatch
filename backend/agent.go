package procwatch

import (
	"context"
	"time"
)

// Agent is an abstract agent, i.e. a service running on a machine and sending
// process lists to a server.
type Agent struct {
	// Connector opens a connection to the server. Must not be nil.
	Connector      Connector

	// ProcessList fetches the process list. Must not be nil.
	ProcessList    ProcessList

	// HostIDProvider returns the host's ID. Must not be nil.
	HostIDProvider HostIDProvider

	// ErrorHandler handles errors. If nil, errors are ignored (in some cases).
	ErrorHandler   func(error)
}

// Run starts the agent.
func (agent *Agent) Run() error {
	hostID := ""
	var err error
	for {
		hostID, err = agent.HostIDProvider.HostID()
		if err != nil {
			agent.handleError(err)
			time.Sleep(sleepInterval)
		} else {
			break
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

		err = client.Identify(context.TODO(), DeviceID(hostID))
		if err != nil {
			agent.handleError(err)
			time.Sleep(sleepInterval)
			continue
		}

		for {
			processNamesFilter, err := client.ProcessNamesFilter(context.TODO())
			if err != nil {
				agent.handleError(err)
				break
			}

			processes, err := agent.ProcessList.Current()
			if err != nil {
				agent.handleError(err)
				break
			}
			agent.handleError(client.Processes(context.TODO(), processes.Filtered(processNamesFilter)))
			time.Sleep(sleepInterval)
		}
	}
}

func (agent *Agent) handleError(err error) {
	if err == nil {
		return
	}
	if handle := agent.ErrorHandler; handle != nil {
		handle(err)
	}
}

const sleepInterval = time.Second

type ProcessList interface {
	// Current returns the processes currently running.
	Current() (Processes, error)
}

// ProcessListFunc implements ProcessList via function.
type ProcessListFunc func() (Processes, error)

// Current calls f and returns the result.
func (f ProcessListFunc) Current() (Processes, error) {
	return f()
}

type HostIDProvider interface {
	// HostID returns the host's ID.
	HostID() (string, error)
}

// HostIDProviderFunc implements HostIDProvider via function.
type HostIDProviderFunc func() (string, error)

// HostID calls f and returns the result.
func (f HostIDProviderFunc) HostID() (string, error) {
	return f()
}

type Connector interface {
	// Connect connects to the server and returns a client.
	Connect() (AgentClient, error)
}

// AgentClient communicates with the server.
type AgentClient interface {
	// Close closes the connection to the server.
	Close() error

	// Identify sends an identification to the server.
	Identify(ctx context.Context, id DeviceID) error

	// ProcessNamesFilter fetches the process names filter from the server. Only processes
	// with names from the list are sent to the server.
	ProcessNamesFilter(ctx context.Context) (ProcessFilterNameList, error)

	// Processes sends a list of processes to the server.
	Processes(ctx context.Context, procs Processes) error
}
