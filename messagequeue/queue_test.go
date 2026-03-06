package messagequeue

import (
	"sync"
	"testing"
)

func TestQueueBasic(t *testing.T) {
	q := New[int]()
	if ok := q.Enqueue(1); !ok {
		t.Fatalf("enqueue failed")
	}

	v, ok := q.Dequeue()
	if !ok || v != 1 {
		t.Fatalf("expected dequeue value=1 ok=true, got value=%d ok=%v", v, ok)
	}
}

func TestQueueCloseAndDrain(t *testing.T) {
	q := New[int]()
	q.Enqueue(10)
	q.Enqueue(20)
	q.Close()

	if _, ok := q.Dequeue(); !ok {
		t.Fatalf("expected first value")
	}
	if _, ok := q.Dequeue(); !ok {
		t.Fatalf("expected second value")
	}
	if _, ok := q.Dequeue(); ok {
		t.Fatalf("expected queue drained and closed")
	}
	if ok := q.Enqueue(30); ok {
		t.Fatalf("expected enqueue on closed queue to fail")
	}
}

func TestQueueConcurrentProducersConsumers(t *testing.T) {
	const producers = 4
	const perProducer = 250
	const total = producers * perProducer

	q := New[int]()

	var producerWG sync.WaitGroup
	for p := 0; p < producers; p++ {
		producerWG.Add(1)
		go func(base int) {
			defer producerWG.Done()
			for i := 0; i < perProducer; i++ {
				q.Enqueue(base*10000 + i)
			}
		}(p)
	}

	results := make(chan int, total)
	var consumerWG sync.WaitGroup
	consumerWG.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer consumerWG.Done()
			for {
				v, ok := q.Dequeue()
				if !ok {
					return
				}
				results <- v
			}
		}()
	}

	producerWG.Wait()
	q.Close()
	consumerWG.Wait()
	close(results)

	seen := map[int]struct{}{}
	for v := range results {
		seen[v] = struct{}{}
	}
	if len(seen) != total {
		t.Fatalf("expected %d unique values, got %d", total, len(seen))
	}
}
