package core

import "syscall"

type FDComm struct {
	Fd int
}

func (f FDComm) Read(data []byte) (int, error) {
	return syscall.Read(f.Fd, data)
}

func (f FDComm) Write(data []byte) (int, error) {
	return syscall.Write(f.Fd, data)
}
