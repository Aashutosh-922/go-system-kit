package eventbus

type Event interface{}

type Handler func(Event)

type Bus struct {
	subscribers map[string][]Handler
}

func New() *Bus {
	return &Bus{
		subscribers: make(map[string][]Handler),
	}
}

func (b *Bus) Subscribe(topic string, h Handler) {
	b.subscribers[topic] = append(b.subscribers[topic], h)
}

func (b *Bus) Publish(topic string, e Event) {
	for _, h := range b.subscribers[topic] {
		go h(e)
	}
}

//Usage

// bus.Subscribe("trade", func(e Event) {
//     log.Println("trade received")
// })

// bus.Publish("trade", tradeEvent)
