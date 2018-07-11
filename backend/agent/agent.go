package agent

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/rpc"

	"net"
	"time"
)

// Agent runs procwatch as an agent.
type Agent struct {
	ServerRPCAddr  string
	ProcessList    ProcessList
	HostIDProvider HostIDProvider
}

// Run connects to the server's address via RPC, registers its ID and sends process informations.
func (ag *Agent) Run() error {
	return (&agent{
		Connector: &RPCConnector{
			ServerRPCAddr: ag.ServerRPCAddr,
		},
		ProcessList:    ag.ProcessList,
		HostIDProvider: ag.HostIDProvider,
	}).Run()
}

type agent struct {
	Connector      Connector
	ProcessList    ProcessList
	HostIDProvider HostIDProvider
}

func (agent *agent) Run() error {
	hostID, err := agent.HostIDProvider.HostID()
	if err != nil {
		return err
	}

	rpcClient, err := agent.Connector.Connect()
	if err != nil {
		return err
	}
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

type HostIDProvider interface {
	HostID() (string, error)
}

type HostIDProviderFunc func() (string, error)

func (f HostIDProviderFunc) HostID() (string, error) {
	return f()
}

type Connector interface {
	Connect() (*rpc.Client, error)
}

type RPCConnector struct {
	ServerRPCAddr string
}

func (connector *RPCConnector) Connect() (*rpc.Client, error) {
	conn, err := net.Dial("tcp", connector.ServerRPCAddr)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn), nil
}
