package messagequeue

import (
	"container/list"
	"sync"
)

// Queue is a concurrency-safe FIFO queue with blocking Dequeue.
type Queue[T any] struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	items    *list.List
	closed   bool
}

func New[T any]() *Queue[T] {
	q := &Queue[T]{
		items: list.New(),
	}
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

// Enqueue adds a new item and returns false if the queue is closed.
func (q *Queue[T]) Enqueue(item T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return false
	}

	q.items.PushBack(item)
	q.notEmpty.Signal()
	return true
}

// Dequeue blocks until an item is available or queue is closed.
// It returns ok=false only when the queue is closed and drained.
func (q *Queue[T]) Dequeue() (item T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.items.Len() == 0 && !q.closed {
		q.notEmpty.Wait()
	}

	if q.items.Len() == 0 && q.closed {
		var zero T
		return zero, false
	}

	front := q.items.Front()
	item = front.Value.(T)
	q.items.Remove(front)
	return item, true
}

func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}
	q.closed = true
	q.notEmpty.Broadcast()
}

func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.items.Len()
}

func (q *Queue[T]) Closed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}

// Usage
//
// q := messagequeue.New[string]()
// go func() {
//     q.Enqueue("job-1")
//     q.Close()
// }()
//
// for {
//     msg, ok := q.Dequeue()
//     if !ok {
//         break
//     }
//     fmt.Println("process:", msg)
// }
