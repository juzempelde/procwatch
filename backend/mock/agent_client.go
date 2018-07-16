package mock

import (
	procwatch "github.com/juzempelde/procwatch/backend"

	"context"
)

type AgentClient struct {
	CloseFunc              func() error
	IdentifyFunc           func(ctx context.Context, id procwatch.DeviceID) error
	ProcessNamesFilterFunc func(ctx context.Context) (procwatch.ProcessFilterNameList, error)
	ProcessesFunc          func(ctx context.Context, procs procwatch.Processes) error
}

func (client *AgentClient) Close() error {
	return client.CloseFunc()
}

func (client *AgentClient) Identify(ctx context.Context, id procwatch.DeviceID) error {
	return client.IdentifyFunc(ctx, id)
}

func (client *AgentClient) ProcessNamesFilter(ctx context.Context) (procwatch.ProcessFilterNameList, error) {
	return client.ProcessNamesFilterFunc(ctx)
}

func (client *AgentClient) Processes(ctx context.Context, procs procwatch.Processes) error {
	return client.ProcessesFunc(ctx, procs)
}
