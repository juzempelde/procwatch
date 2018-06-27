package wrap

import (
	"github.com/juzempelde/procwatch/backend/agent"
	"github.com/juzempelde/procwatch/backend/cli"
)

// Agent wraps an agent into something the CLI package can use.
func Agent(ag *agent.Agent) cli.Agent {
	return &cliAgent{
		agent: ag,
	}
}

type cliAgent struct {
	agent *agent.Agent
}

func (ag *cliAgent) Exec() error {
	return ag.agent.Run()
}

func (ag *cliAgent) SetServerRPCAddr(addr string) {
	ag.agent.ServerRPCAddr = addr
}
