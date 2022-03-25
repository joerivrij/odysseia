package queue

import (
	"context"
	"github.com/kubemq-io/kubemq-go"
)

func (k *KubeMq) Close() error {
	return k.client.Close()
}

func (k *KubeMq) SendQueueMessage(body []byte) (*kubemq.SendQueueMessageResult, error) {
	message := kubemq.NewQueueMessage().
		SetChannel(k.channel).SetBody(body)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := k.client.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (k *KubeMq) PullMessages(numberOfMessages int32, id string) (*kubemq.ReceiveQueueMessagesResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messages, err := k.client.Pull(ctx, &kubemq.ReceiveQueueMessagesRequest{
		ClientID:            id,
		Channel:             k.channel,
		MaxNumberOfMessages: numberOfMessages,
		WaitTimeSeconds:     1,
	})

	return messages, err
}

func (k *KubeMq) PullRemainingMessages(id string) (*kubemq.ReceiveQueueMessagesResponse, error) {
	info, err := k.Info()
	if err != nil {
		return nil, err
	}

	messagesRemaining := info.Waiting

	messages, err := k.PullMessages(int32(messagesRemaining), id)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (k *KubeMq) Info() (*kubemq.QueuesInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := k.client.QueuesInfo(ctx, k.channel)
	if err != nil {
		return nil, err
	}

	return res, nil
}
