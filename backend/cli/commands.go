package cli

// NewCommands wraps functions for creating Agent and Server into a Commands instance.
func NewCommands(
	agentFunc func() Agent,
	serverFunc func() Server,
) Commands {
	return &commands{
		agentFunc:  agentFunc,
		serverFunc: serverFunc,
	}
}

type commands struct {
	agentFunc  func() Agent
	serverFunc func() Server
}

func (cmd *commands) Agent() Agent {
	return cmd.agentFunc()
}

func (cmd *commands) Server() Server {
	return cmd.serverFunc()
}
