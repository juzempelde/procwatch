package agent

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/rpc"

	"net"
	"os"
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

	_, err = rpcClient.ProcessNames()
	if err != nil {
		return err
	}

	return nil
}
