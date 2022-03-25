package queue

import (
	"context"
	"github.com/kubemq-io/kubemq-go"
)

type Client interface {
	Close() error
	SendQueueMessage(body []byte) (*kubemq.SendQueueMessageResult, error)
	PullMessages(numberOfMessages int32, id string) (*kubemq.ReceiveQueueMessagesResponse, error)
	PullRemainingMessages(id string) (*kubemq.ReceiveQueueMessagesResponse, error)
	Info() (*kubemq.QueuesInfo, error)
}

type KubeMq struct {
	client  *kubemq.QueuesClient
	channel string
}

type FakeKubeMq struct {
}

func NewKubeMqClient(address, id, channel string, port int) (Client, error) {
	queuesClient, err := kubemq.NewQueuesStreamClient(context.Background(),
		kubemq.WithAddress(address, port),
		kubemq.WithClientId(id),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		return nil, err
	}

	kubeMq := KubeMq{
		client:  queuesClient,
		channel: channel,
	}

	return &kubeMq, nil
}

func NewFakeKubeMqClient() (Client, error) {
	return &FakeKubeMq{}, nil
}
