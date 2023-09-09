package core

import (
	"log"
	"memkv/internal/config"
)

type Epoll struct {
	fd          int
	epollEvents []syscall.EpollEvent
}

func Create() (*Epoll, error) {
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Epoll{
		fd:          epollFD,
		epollEvents: make([]syscall.EpollEvent, config.MaxConnection),
	}, nil
}

func (ep *Epoll) Monitor(event Event) error {
	epollEvent := event.toNative()
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, event.Fd, &epollEvent)
}
