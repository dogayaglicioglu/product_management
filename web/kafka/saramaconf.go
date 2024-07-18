package kafka

import "github.com/IBM/sarama"

func SaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Return.Errors = true
	return config
}
