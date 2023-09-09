package core

import "syscall"

func (e Event) toNative() syscall.EpollEvent {
	return syscall.EpollEvent{
		Fd:     int32(e.Fd),
		Events: uint32(e.Op),
	}
}

func createEvent(ep syscall.EpollEvent) Event {
	return Event{
		Fd: int(ep.Fd),
		Op: int(ep.Events),
	}
}
