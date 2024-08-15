package usecase

import (
	"encoding/json"
	"github.com/cenkalti/backoff/v4"
	"log"
	"orchestra-svc/internal/dto/event"
	"time"
)

func (o *OrchestraUsecase) retryKafkaSend(topic, key string, value []byte) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Minute

	return backoff.RetryNotify(func() error {
		log.Println("Sending message to Kafka")
		err := o.producer.SendMessage(topic, key, value)
		if err != nil {

			var gevent event.GlobalEvent[any, any]
			if jsonErr := json.Unmarshal(value, &gevent); jsonErr == nil {
				if gevent.StatusCode == 500 {
					return err
				}
			}
			return backoff.Permanent(err)
		}
		return nil
	}, b, func(err error, duration time.Duration) {
		log.Printf("Retrying Kafka send in %v due to error: %v", duration, err)
	})
}
