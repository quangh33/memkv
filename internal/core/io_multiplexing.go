package core

type Event struct {
	Fd int
	Op int
}

type ioMultiplexing interface {
	Monitor(event Event) error
	Wait() ([]Event, error)
	Close() error
}
