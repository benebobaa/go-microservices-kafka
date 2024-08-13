package kafka

import (
	"context"
	"log"
	"product-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  sarama.ConsumerGroupHandler
}

type MessageHandler struct {
	userUsecase *usecase.Usecase
	producer    *KafkaProducer
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		log.Printf("Received message: %s", string(msg.Value))

		sess.MarkMessage(msg, "")
	}
	return nil
}

func NewKafkaConsumer(
	brokers []string, groupID string,
	topics []string, usecase *usecase.Usecase,
	producer *KafkaProducer,
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
		handler: MessageHandler{
			userUsecase: usecase,
			producer:    producer,
		},
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
