package procwatch

import (
	"context"
	"time"
)

type agentClientWithDeadline struct {
	client AgentClient

	Timeout time.Duration
}

// AgentClientWithDeadline wraps a client and adds timeouts. Useful for clients which do not implement timeouts themselves.
func AgentClientWithDeadline(client AgentClient, timeout time.Duration) AgentClient {
	return &agentClientWithDeadline{
		client:  client,
		Timeout: timeout,
	}
}

func (client *agentClientWithDeadline) Close() error {
	return client.client.Close()
}

func (client *agentClientWithDeadline) Identify(ctx context.Context, id DeviceID) error {
	return client.timeout(
		ctx,
		func() error {
			return client.client.Identify(ctx, id)
		},
	)
}

func (client *agentClientWithDeadline) ProcessNamesFilter(ctx context.Context) (ProcessFilterNameList, error) {
	var list ProcessFilterNameList
	err := client.timeout(
		ctx,
		func() error {
			var opErr error
			list, opErr = client.client.ProcessNamesFilter(ctx)
			return opErr
		},
	)
	return list, err
}

func (client *agentClientWithDeadline) Processes(ctx context.Context, procs Processes) error {
	return client.timeout(
		ctx,
		func() error {
			return client.client.Processes(ctx, procs)
		},
	)
}

func (client *agentClientWithDeadline) timeout(ctx context.Context, op func() error) error {
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(client.Timeout))
	defer cancelFunc()
	opFinished := make(chan struct{})
	var opErr error
	go func() {
		opErr = op()
		close(opFinished)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-opFinished:
		return opErr
	}
}
