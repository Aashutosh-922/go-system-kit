package main

import (
	"fmt"
	"time"

	"github.com/Aashutosh-922/go-system-kit.git/messagequeue"
)

func main() {
	q := messagequeue.New[string]()

	go func() {
		for i := 1; i <= 5; i++ {
			job := fmt.Sprintf("job-%d", i)
			if ok := q.Enqueue(job); !ok {
				return
			}
			time.Sleep(150 * time.Millisecond)
		}
		q.Close()
	}()

	for {
		msg, ok := q.Dequeue()
		if !ok {
			break
		}
		fmt.Println("processing", msg)
	}

	fmt.Println("queue closed and drained")
}
