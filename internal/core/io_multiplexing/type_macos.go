//go:build darwin

package io_multiplexing

import "syscall"

func (e Event) toNative(flags uint16) syscall.Kevent_t {
	var filter int16 = syscall.EVFILT_WRITE
	if e.Op == OpRead {
		filter = syscall.EVFILT_READ
	}
	return syscall.Kevent_t{
		Ident:  uint64(e.Fd),
		Filter: filter,
		Flags:  flags,
	}
}

func createEvent(kq syscall.Kevent_t) Event {
	var op Operation = OpWrite
	if kq.Filter == syscall.EVFILT_READ {
		op = OpRead
	}
	return Event{
		Fd: int(kq.Ident),
		Op: op,
	}
}
