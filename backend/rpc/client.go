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

func (client *Client) Identify(ctx context.Context, id procwatch.DeviceID) error {
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

func (client *Client) ProcessNamesFilter(ctx context.Context) (procwatch.ProcessFilterNameList, error) {
	request := ProcessNameFilterRequest{}
	response := &ProcessNameFilterResponse{}
	err := client.caller.Call(fmt.Sprintf("%s.%s", processFilter, "Expose"), request, response)
	if err != nil {
		return nil, err
	}
	return procwatch.ProcessFilterNameList(response.Names), nil
}

func (client *Client) Processes(ctx context.Context, procs procwatch.Processes) error {
	request := ProcessesRequest{
		Processes: procs,
	}
	response := &ProcessesResponse{}
	return client.caller.Call(fmt.Sprintf("%s.%s", processes, "Processes"), request, response)
}

func (client *Client) Close() error {
	return client.caller.Close()
}
