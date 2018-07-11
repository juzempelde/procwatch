package rpc

import (
	"github.com/juzempelde/procwatch/backend"

	"fmt"
	"io"
	"net/rpc"
)

// caller abstracts away a net/rpc client.
type caller interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
	Close() error
}

type Client struct {
	caller caller
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{
		caller: rpc.NewClient(conn),
	}
}

func (client *Client) Identify(id procwatch.DeviceID) error {
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

func (client *Client) ProcessNamesFilter() (procwatch.ProcessFilterNameList, error) {
	request := ProcessNameFilterRequest{}
	response := &ProcessNameFilterResponse{}
	err := client.caller.Call(fmt.Sprintf("%s.%s", processFilter, "ProcessNames"), request, response)
	if err != nil {
		return nil, err
	}
	return procwatch.ProcessFilterNameList(response.Names), nil
}

func (client *Client) Processes(procs procwatch.Processes) error {
	request := ProcessesRequest{
		Processes: procs,
	}
	response := &ProcessesResponse{}
	return client.caller.Call(fmt.Sprintf("%s.%s", processes, "Processes"), request, response)
}

func (client *Client) Close() error {
	return client.caller.Close()
}
