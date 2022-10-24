package aristoteles

import (
	"github.com/odysseia-greek/plato/queue"
	"strconv"
)

func (c *Config) getMqQueueClient() (queue.Client, error) {
	channel := c.getStringFromEnv(EnvChannel, defaultChannelName)
	address := c.getStringFromEnv(EnvMqAddress, defaultMqAddress)
	portFromEnv := c.getStringFromEnv(EnvMqPort, defaultMqPort)
	id := c.getStringFromEnv(EnvPodName, defaultPodName)
	port, err := strconv.Atoi(portFromEnv)
	if err != nil {
		return nil, err
	}

	kubeMq, err := queue.NewKubeMqClient(address, id, channel, port)
	if err != nil {
		return nil, err
	}

	return kubeMq, nil
}
