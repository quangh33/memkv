package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	"memkv/internal/config"
	"memkv/internal/constant"
	"memkv/internal/core"
)

var eStatus int32 = constant.EngineStatus_WAITING

func WaitForSignal(wg *sync.WaitGroup, signals chan os.Signal) {
	defer wg.Done()
	<-signals
	for atomic.LoadInt32(&eStatus) == constant.EngineStatus_BUSY {
	}

	if !atomic.CompareAndSwapInt32(&eStatus, constant.EngineStatus_WAITING, constant.EngineStatus_SHUTTING_DOWN) {
		// rarely happen
		log.Println("Engine is busy again. Try again!")
		return
	}
	core.Shutdown()
	os.Exit(0)
}

func readCommandFD(fd int) (*core.MemKVCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	return core.ParseCmd(buf[:n])
}

func responseRw(cmd *core.MemKVCmd, rw io.ReadWriter) {
	err := core.EvalAndResponse(cmd, rw)
	if err != nil {
		responseErrorRw(err, rw)
	}
}

func responseErrorRw(err error, rw io.ReadWriter) {
	rw.Write([]byte(fmt.Sprintf("-%s%s", err, core.CRLF)))
}

func RunAsyncTCPServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Println("starting an asynchronous TCP server on", config.Host, config.Port)

	// Create EPOLL Event Objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, config.MaxConnection)
	client_number := 0

	// Create a server socket - an endpoint for communication between client and server
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Println(err)
		return err
	}
	defer syscall.Close(serverFD)

	// Set the Socket operate in a non-blocking mode
	// Default mode is blocking mode: when you read from a FD, control isn't returned
	// until at least one byte of data is read.
	// Non-blocking mode: if the read buffer is empty, it will return immediately.
	// We want non-blocking mode because we will use epoll to monitor and then read from
	// multiple FD, so we want to ensure that none of them cause the program to "lock up."
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		log.Println(err)
		return err
	}

	// Bind the IP and the port to the server socket FD.
	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		log.Println(err)
		return err
	}

	// Start listening
	if err = syscall.Listen(serverFD, config.MaxConnection); err != nil {
		log.Println(err)
		return err
	}

	// creating EPOLL instance
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epollFD)

	// Specify the events we want to monitor on server socket FD
	// here we want to get hint whenever server socket FD is available for read operation
	var socketServerReadReadyEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	// Listen to read events on the Server itself
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerReadReadyEvent); err != nil {
		return err
	}

	for atomic.LoadInt32(&eStatus) != constant.EngineStatus_SHUTTING_DOWN {
		// see if any FD is ready for an IO
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		if !atomic.CompareAndSwapInt32(&eStatus, constant.EngineStatus_WAITING, constant.EngineStatus_BUSY) {
			if eStatus == constant.EngineStatus_SHUTTING_DOWN {
				return nil
			}
		}
		for i := 0; i < nevents; i++ {
			// if the socket server itself is ready for an IO
			if int(events[i].Fd) == serverFD {
				// accept the incoming connection from a client
				client_number++
				log.Printf("new client: id=%d\n", client_number)
				connFD, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				if err = syscall.SetNonblock(connFD, true); err != nil {
					return err
				}

				// add this new TCP connection to be monitored
				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(connFD),
				}
				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, connFD, &socketClientEvent); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommandFD(comm.Fd)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					client_number--
					log.Println("client quit")
					continue
				}
				responseRw(cmd, comm)
			}
			atomic.SwapInt32(&eStatus, constant.EngineStatus_WAITING)
		}
	}

	return nil
}
