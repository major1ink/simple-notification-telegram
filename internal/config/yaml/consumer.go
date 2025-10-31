package yaml

import (
	"github.com/IBM/sarama"
)

type ConsumerConfig struct {
	Topic   string `yaml:"topic"`
	GroupId string `yaml:"group_id"`
}

func (c *ConsumerConfig) GetTopic() string {
	return c.Topic
}

func (c *ConsumerConfig) GetGroupId() string {
	return c.GroupId
}

func (с *ConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
