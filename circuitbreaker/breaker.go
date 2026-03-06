package circuit

import "time"

type Breaker struct {
	failures int
	open     bool
}

func (c *Breaker) Call(fn func() error) error {

	if c.open {
		return nil
	}

	err := fn()

	if err != nil {
		c.failures++

		if c.failures > 5 {
			c.open = true

			go func() {
				time.Sleep(5 * time.Second)
				c.open = false
				c.failures = 0
			}()
		}
	}

	return err
}
