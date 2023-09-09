package core

type Event struct {
	Fd int
	Op int
}

type IOMultiplexer interface {
	Monitor(event Event) error
	Check() ([]Event, error)
	Close() error
}
