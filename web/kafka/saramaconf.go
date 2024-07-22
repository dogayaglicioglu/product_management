package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

func SaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Group.Session.Timeout = 30 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 10 * time.Second
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_1_0_0
	return config
}
