package agent

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/rpc"

	"fmt"
	"net"
)

// Agent runs procwatch as an agent.
type Agent struct {
	ServerRPCAddr  string
	ProcessList    procwatch.ProcessList
	HostIDProvider procwatch.HostIDProvider
}

// Run connects to the server's address via RPC, registers its ID and sends process informations.
func (ag *Agent) Run() error {
	return (&procwatch.Agent{
		Connector: &RPCConnector{
			ServerRPCAddr: ag.ServerRPCAddr,
		},
		ProcessList:    ag.ProcessList,
		HostIDProvider: ag.HostIDProvider,
		ErrorHandler: func(err error) {
			fmt.Printf("Error: %+v\n", err)
		},
	}).Run()
}

type RPCConnector struct {
	ServerRPCAddr string
}

func (connector *RPCConnector) Connect() (procwatch.AgentClient, error) {
	conn, err := net.Dial("tcp", connector.ServerRPCAddr)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn), nil
}
