package agent

import (
	"net"
	"net/rpc"
)

// Agent runs procwatch as an agent.
type Agent struct {
	ServerRPCAddr string
}

// Run connects to the server's address via RPC, registers its ID and sends process informations.
func (agent *Agent) Run() error {
	conn, err := net.Dial("tcp", agent.ServerRPCAddr)
	if err != nil {
		return err
	}
	rpcClient := rpc.NewClient(conn)
	_ = rpcClient.Close() // TODO: Handle close error
	return nil
}
