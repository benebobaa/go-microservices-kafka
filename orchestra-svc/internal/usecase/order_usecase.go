package usecase

import "orchestra-svc/pkg/producer"

type OrderUsecase struct {
	producer *producer.KafkaProducer
}

func NewOrderUsecase(producer *producer.KafkaProducer) *OrderUsecase {
	return &OrderUsecase{
		producer: producer,
	}
}
