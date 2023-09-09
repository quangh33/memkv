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

var eStatus int32 = constant.EngineStatusWaiting

func WaitForSignal(wg *sync.WaitGroup, signals chan os.Signal) {
	defer wg.Done()
	<-signals
	for atomic.LoadInt32(&eStatus) == constant.EngineStatusBusy {
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

	var events = make([]core.Event, config.MaxConnection)
	clientNumber := 0

	// Create a server socket. A socket is an endpoint for communication between client and server
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
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

	// ioMultiplexer is an object that can monitor multiple file descriptor (FD) at the same time.
	// When one or more monitored FD(s) are ready for IO, it will notify our server.
	// Here, we use ioMultiplexer to monitor Server FD and Clients FD.
	ioMultiplexer, err := core.CreateIOMultiplexer()
	if err != nil {
		return err
	}
	defer ioMultiplexer.Close()

	// Monitor "read" events on the Server FD
	if err = ioMultiplexer.Monitor(core.Event{
		Fd: serverFD,
		Op: core.OpRead,
	}); err != nil {
		return err
	}

	for atomic.LoadInt32(&eStatus) != constant.EngineStatusShuttingDown {
		// check if any FD is ready for an IO
		events, err = ioMultiplexer.Check()
		if err != nil {
			continue
		}

		if !atomic.CompareAndSwapInt32(&eStatus, constant.EngineStatusWaiting, constant.EngineStatusBusy) {
			if eStatus == constant.EngineStatusShuttingDown {
				return nil
			}
		}
		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFD {
				// the Server FD is ready for an IO, means we have a new client.
				clientNumber++
				log.Printf("new client: id=%d\n", clientNumber)
				// accept the incoming connection from a client
				connFD, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				if err = syscall.SetNonblock(connFD, true); err != nil {
					return err
				}

				// add this new connection to be monitored
				if err = ioMultiplexer.Monitor(core.Event{
					Fd: connFD,
					Op: core.OpRead,
				}); err != nil {
					return err
				}
			} else {
				// the Client FD is ready for an IO, means an existing client is sending a command
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommandFD(comm.Fd)
				if err != nil {
					syscall.Close(events[i].Fd)
					clientNumber--
					log.Println("client quit")
					atomic.SwapInt32(&eStatus, constant.EngineStatusWaiting)
					continue
				}
				responseRw(cmd, comm)
			}
			atomic.SwapInt32(&eStatus, constant.EngineStatusWaiting)
		}
	}

	return nil
}
