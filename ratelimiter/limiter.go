package ratelimiter

import "time"

type Limiter struct {
	tokens chan struct{}
}

func New(rate int) *Limiter {

	l := &Limiter{
		tokens: make(chan struct{}, rate),
	}

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rate))

		for range ticker.C {
			select {
			case l.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return l
}

func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

//Usage

// if !limiter.Allow() {
//     http.Error(w, "rate limit exceeded", 429)
// }
