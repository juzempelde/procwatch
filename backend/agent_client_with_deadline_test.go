package procwatch_test

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/mock"

	"context"
	"fmt"
	"testing"
	"time"
)

func TestAgentClientWithDeadlinePassesClose(t *testing.T) {
	client := procwatch.AgentClientWithDeadline(
		&mock.AgentClient{
			CloseFunc: func() error {
				return fmt.Errorf("close failed")
			},
		},
		time.Microsecond,
	)

	err := client.Close()

	if err == nil {
		t.Errorf("Expected error not to be nil")
	}
}

func TestAgentClientWithDeadlinePassesProcessNamesFilter(t *testing.T) {
	returnedList := procwatch.ProcessFilterNameList{"foo", "bar", "baz"}
	client := procwatch.AgentClientWithDeadline(
		&mock.AgentClient{
			ProcessNamesFilterFunc: func(ctx context.Context) (procwatch.ProcessFilterNameList, error) {
				return returnedList, nil
			},
		},
		time.Microsecond,
	)

	list, err := client.ProcessNamesFilter(context.Background())

	if len(list) != 3 {
		t.Errorf("Expected list to consist of %+v, but got %+v", returnedList, list)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %+v", err)
	}
}

func TestAgentClientWithDeadlinePassesIdentify(t *testing.T) {
	var receivedID procwatch.DeviceID
	client := procwatch.AgentClientWithDeadline(
		&mock.AgentClient{
			IdentifyFunc: func(ctx context.Context, id procwatch.DeviceID) error {
				receivedID = id
				return fmt.Errorf("setting id failed")
			},
		},
		time.Microsecond,
	)

	inputID := procwatch.DeviceID("xyz")
	err := client.Identify(context.Background(), inputID)

	if err == nil {
		t.Errorf("Expected error not to be nil")
	}
	if receivedID != inputID {
		t.Errorf("Expected ID to be %s, but got %s", inputID, receivedID)
	}
}

func TestAgentClientWithDeadlineCancelsIdentify(t *testing.T) {
	var receivedID procwatch.DeviceID
	finishIdentify := make(chan struct{})
	client := procwatch.AgentClientWithDeadline(
		&mock.AgentClient{
			IdentifyFunc: func(ctx context.Context, id procwatch.DeviceID) error {
				<-finishIdentify
				receivedID = id
				return nil
			},
		},
		time.Microsecond,
	)

	inputID := procwatch.DeviceID("someid")
	err := client.Identify(context.Background(), inputID)
	if receivedID != procwatch.DeviceID("") {
		t.Errorf("Expected device ID to be empty, but got %s", receivedID)
	}
	if err == nil {
		t.Errorf("Expected error not be nil")
	}
}
