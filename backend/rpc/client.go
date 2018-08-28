package rpc

import (
	"github.com/juzempelde/procwatch/backend"

	"context"
	"fmt"
	"io"
	"net/rpc"
	"time"
)

// caller abstracts away a net/rpc client.
type caller interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
	Close() error
}

// Client communicates with a procwatch RPC server. The zero value is not useful.
// Must be created with NewClient().
type Client struct {
	caller           caller
	deadlineAcceptor deadlineAcceptor
}

type deadlineAcceptor interface {
	SetDeadline(t time.Time) error
}

type nopDeadlineAcceptor struct{}

func (acc nopDeadlineAcceptor) SetDeadline(t time.Time) error {
	return nil
}

// NewClient creates a new client from an existing ReadWriteCloser.
func NewClient(conn io.ReadWriteCloser) *Client {
	client := &Client{
		caller:           rpc.NewClient(conn),
		deadlineAcceptor: nopDeadlineAcceptor{},
	}
	if da, ok := conn.(deadlineAcceptor); ok {
		client.deadlineAcceptor = da
	}
	return client
}

func (client *Client) setDeadlineFromContext(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil
	}
	return client.deadlineAcceptor.SetDeadline(deadline)
}

// Identify sends an ID to the server which identifies this agent.
func (client *Client) Identify(ctx context.Context, id procwatch.DeviceID) error {
	client.setDeadlineFromContext(ctx)
	request := &IdentificationRequest{
		ID: id,
	}
	response := &IdentificationResponse{}
	err := client.caller.Call(fmt.Sprintf("%s.%s", identificator, "Identify"), request, response)
	if err != nil {
		return err
	}
	if response.ErrReason != "" {
		return fmt.Errorf(response.ErrReason)
	}
	return nil
}

// ProcessNamesFilter retrieves the filter for process names from the server.
func (client *Client) ProcessNamesFilter(ctx context.Context) (procwatch.ProcessFilterNameList, error) {
	client.setDeadlineFromContext(ctx)
	request := ProcessNameFilterRequest{}
	response := &ProcessNameFilterResponse{}
	err := client.caller.Call(fmt.Sprintf("%s.%s", processFilter, "Expose"), request, response)
	if err != nil {
		return nil, err
	}
	return procwatch.ProcessFilterNameList(response.Names), nil
}

// Processes sends a list of processes to the server.
func (client *Client) Processes(ctx context.Context, procs procwatch.Processes) error {
	client.setDeadlineFromContext(ctx)
	request := ProcessesRequest{
		Processes: procs,
	}
	response := &ProcessesResponse{}
	return client.caller.Call(fmt.Sprintf("%s.%s", processes, "Processes"), request, response)
}

// Close closes the connection to the RPC server.
func (client *Client) Close() error {
	return client.caller.Close()
}
