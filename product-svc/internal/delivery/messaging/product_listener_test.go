package messaging

import (
	"context"
	"github.com/IBM/sarama"
	"testing"
)

type mockConsumerGroupSession struct{}

func (m *mockConsumerGroupSession) Commit() {
}

func (m *mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {}
func (m *mockConsumerGroupSession) Context() context.Context                                 { return context.TODO() }
func (m *mockConsumerGroupSession) Claims() map[string][]int32                               { return nil }
func (m *mockConsumerGroupSession) MemberID() string                                         { return "" }
func (m *mockConsumerGroupSession) GenerationID() int32                                      { return 0 }
func (m *mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}
func (m *mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

type mockConsumerGroupClaim struct{}

func (m *mockConsumerGroupClaim) Topic() string              { return "mockTopic" }
func (m *mockConsumerGroupClaim) Partition() int32           { return 0 }
func (m *mockConsumerGroupClaim) InitialOffset() int64       { return 0 }
func (m *mockConsumerGroupClaim) HighWaterMarkOffset() int64 { return 0 }
func (m *mockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	msgChan := make(chan *sarama.ConsumerMessage, 1)
	msgChan <- &sarama.ConsumerMessage{
		Value: []byte("test message"),
	}
	close(msgChan)
	return msgChan
}

func TestConsumeClaim(t *testing.T) {
	consumer := MessageHandler{}
	session := &mockConsumerGroupSession{}
	claim := &mockConsumerGroupClaim{}

	err := consumer.ConsumeClaim(session, claim)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
