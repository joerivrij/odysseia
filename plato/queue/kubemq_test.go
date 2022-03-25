package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMqClient(t *testing.T) {
	address := "localhost"
	port := 50000
	id := "testid"
	channel := "testchannel"
	message := "this is a test message"

	body := []byte(fmt.Sprintf("sending message %s", message))

	t.Run("SendMessageWithRealClient", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping testing in short mode")
		}

		testClient, err := NewKubeMqClient(address, id, channel, port)
		assert.Nil(t, err)

		sut, err := testClient.SendQueueMessage(body)
		assert.Nil(t, err)
		assert.False(t, sut.IsError)
	})

	t.Run("Receive", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping testing in short mode")
		}

		testClient, err := NewKubeMqClient(address, id, channel, port)
		assert.Nil(t, err)

		receiveMessages := 1

		sut, err := testClient.PullMessages(int32(receiveMessages), id)
		assert.Nil(t, err)
		assert.NotNil(t, sut)
	})

	t.Run("Info", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping testing in short mode")
		}

		testClient, err := NewKubeMqClient(address, id, channel, port)
		assert.Nil(t, err)

		sut, err := testClient.Info()
		assert.Nil(t, err)
		assert.NotNil(t, sut)
	})

	t.Run("SendMessageWithFakeClient", func(t *testing.T) {
		testClient, err := NewFakeKubeMqClient()
		assert.Nil(t, err)

		sut, err := testClient.SendQueueMessage(body)
		assert.Nil(t, err)
		assert.False(t, sut.IsError)
	})

	t.Run("SendEmptyWithFakeClient", func(t *testing.T) {
		testClient, err := NewFakeKubeMqClient()
		assert.Nil(t, err)

		sut, err := testClient.SendQueueMessage(nil)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})

	t.Run("Receive", func(t *testing.T) {
		testClient, err := NewFakeKubeMqClient()
		assert.Nil(t, err)

		receiveMessages := 1

		sut, err := testClient.PullMessages(int32(receiveMessages), id)
		assert.Nil(t, err)
		assert.NotNil(t, sut)
	})

	t.Run("Info", func(t *testing.T) {
		testClient, err := NewFakeKubeMqClient()
		assert.Nil(t, err)
		sut, err := testClient.Info()
		assert.Nil(t, err)
		assert.NotNil(t, sut)
	})

}
