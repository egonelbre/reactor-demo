package reactor

type ID uint64

type Reactor struct {
	Inbound  Inbound
	Outbound Outbound
}

type Future struct {
	Resolve func(response Message)
}

type Event struct{}

type Message []byte

type Inbound struct {
	Messages []Event
}

type Outbound struct {
	Messages []Event
}

type Reactor interface {
	React(Event)
}

type Event interface{}
