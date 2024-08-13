package consumer

import (
	"context"
	"user-svc/internal/delivery/messaging"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  sarama.ConsumerGroupHandler
}

func NewKafkaConsumer(
	brokers []string, groupID string,
	topics []string, messageHandler *messaging.MessageHandler,
) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topics:   topics,
		handler:  messageHandler,
	}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	for {
		err := kc.consumer.Consume(ctx, kc.topics, kc.handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
