package retry

import (
	"time"
)

func WithBackoff(attempts int, fn func() error) error {

	var err error

	for i := 0; i < attempts; i++ {
		err = fn()

		if err == nil {
			return nil
		}

		time.Sleep(time.Duration(i*i) * time.Millisecond)
	}

	return err
}

// Usage

// retry.WithBackoff(5, func() error {
//     return kafkaProducer.Send(msg)
// })
