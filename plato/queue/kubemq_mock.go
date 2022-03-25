package queue

import (
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"github.com/odysseia/plato/models"
)

func (f *FakeKubeMq) Close() error {
	return nil
}

func (f *FakeKubeMq) SendQueueMessage(body []byte) (*kubemq.SendQueueMessageResult, error) {
	if body == nil {
		return nil, fmt.Errorf("an error was created in the mock")
	}
	res := kubemq.SendQueueMessageResult{
		MessageID:    "yvhAg3YhGz33G5fSErJeF3",
		SentAt:       1646821272326973107,
		ExpirationAt: 0,
		DelayedTo:    0,
		IsError:      false,
		Error:        "",
	}

	return &res, nil
}

func (f *FakeKubeMq) Info() (*kubemq.QueuesInfo, error) {
	queue := kubemq.QueueInfo{
		Name:          "",
		Messages:      0,
		Bytes:         0,
		FirstSequence: 0,
		LastSequence:  0,
		Sent:          0,
		Subscribers:   0,
		Waiting:       0,
		Delivered:     10,
	}

	res := kubemq.QueuesInfo{
		TotalQueues: 0,
		Sent:        0,
		Waiting:     0,
		Delivered:   0,
		Queues: []*kubemq.QueueInfo{
			&queue,
		},
	}

	return &res, nil
}

func (f *FakeKubeMq) PullMessages(numberOfMessages int32, id string) (*kubemq.ReceiveQueueMessagesResponse, error) {
	if numberOfMessages == 0 {
		return nil, fmt.Errorf("a mocking error while receiving messages")
	}

	message := kubemq.NewQueueMessage()
	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "market place",
		LinkedWord: "",
		Original:   "",
	}

	marshalledBody, err := body.Marshal()
	if err != nil {
		return nil, err
	}

	message.Body = marshalledBody
	message.ClientID = id

	messages := kubemq.ReceiveQueueMessagesResponse{
		RequestID:        "",
		Messages:         []*kubemq.QueueMessage{message},
		MessagesReceived: 0,
		MessagesExpired:  0,
		IsPeak:           false,
		IsError:          false,
		Error:            "",
	}

	return &messages, nil
}

func (f *FakeKubeMq) PullRemainingMessages(id string) (*kubemq.ReceiveQueueMessagesResponse, error) {
	message := kubemq.NewQueueMessage()

	message.ClientID = id
	message.Body = []byte("hello from the mock")

	messages := kubemq.ReceiveQueueMessagesResponse{
		RequestID:        "",
		Messages:         []*kubemq.QueueMessage{message},
		MessagesReceived: 0,
		MessagesExpired:  0,
		IsPeak:           false,
		IsError:          false,
		Error:            "",
	}

	return &messages, nil
}
