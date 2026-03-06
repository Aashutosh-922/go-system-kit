package goroutines

import (
    "sync"
)

type Group struct {
    wg sync.WaitGroup
}

func (g *Group) Go(fn func()) {
    g.wg.Add(1)

    go func() {
        defer g.wg.Done()
        fn()
    }()
}

func (g *Group) Wait() {
    g.wg.Wait()
}  //goroutine orchestration helper


//Usage

// var g goroutines.Group

// g.Go(func() {
//     startHTTP()
// })

// g.Go(func() {
//     startKafkaConsumer()
// })

// g.Wait()