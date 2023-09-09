//go:build linux

package io_multiplexing

import "syscall"

func (e Event) toNative() syscall.EpollEvent {
	var event uint32 = syscall.EPOLLIN
	if e.Op == OpWrite {
		event = syscall.EPOLLOUT
	}
	return syscall.EpollEvent{
		Fd:     int32(e.Fd),
		Events: event,
	}
}

func createEvent(ep syscall.EpollEvent) Event {
	var op Operation = OpRead
	if ep.Events == syscall.EPOLLOUT {
		op = OpWrite
	}
	return Event{
		Fd: int(ep.Fd),
		Op: op,
	}
}
