package rpc_test

import (
	"fmt"

	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/mock"
	"github.com/juzempelde/procwatch/backend/rpc"

	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestSuccessfulIdentification(t *testing.T) {
	serverConn, clientConn, err := connect()
	if err != nil {
		t.Fatalf("Could not create connections: %+v", err)
	}
	defer serverConn.Close()
	defer clientConn.Close()
	server := rpc.NewServer(&mock.Device{}, nopRefreshDeadline)
	go server.ServeConn(serverConn)

	time.Sleep(time.Millisecond)
	client := rpc.NewClient(clientConn)

	err = client.Identify(context.TODO(), procwatch.DeviceID("xyz"))

	if err != nil {
		t.Errorf("Expected no error, but got %+v", err)
	}
}

func TestFailingIdentification(t *testing.T) {
	serverConn, clientConn, err := connect()
	if err != nil {
		t.Fatalf("Could not create connections: %+v", err)
	}
	defer serverConn.Close()
	defer clientConn.Close()
	server := rpc.NewServer(&mock.Device{IdentifyErrors: map[procwatch.DeviceID]error{"abc": fmt.Errorf("not possible")}}, nopRefreshDeadline)
	go server.ServeConn(serverConn)

	time.Sleep(time.Millisecond)
	client := rpc.NewClient(clientConn)

	err = client.Identify(context.TODO(), procwatch.DeviceID("abc"))

	if err == nil {
		t.Errorf("Expected error")
	}
}

func connect() (net.Conn, net.Conn, error) {
	var acceptConn, dialConn net.Conn
	var acceptErr, dialErr error

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		acceptConn, acceptErr = listener.Accept()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		dialConn, dialErr = net.Dial(listener.Addr().Network(), listener.Addr().String())
	}()

	wg.Wait()

	if acceptErr != nil {
		return nil, nil, acceptErr
	}
	if dialErr != nil {
		return nil, nil, dialErr
	}

	return acceptConn, dialConn, nil
}

func nopRefreshDeadline() error {
	return nil
}
