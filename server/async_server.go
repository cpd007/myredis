package server

import (
	"log"
	"net"

	"github.com/cpd007/myredis/config"
	"github.com/cpd007/myredis/core"
	"golang.org/x/sys/unix"
)

var con_clients uint32 = 0

func StartAsyncTCPServer() {

	log.Printf("Starting asynchronous tcp server")

	max_clients := 20000

	// create the socket
	serverFd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal("error in creating the socket: ", err)
	}
	defer unix.Close(serverFd)

	// set the scocket to non-blocking mode
	err = unix.SetNonblock(serverFd, true)
	if err != nil {
		log.Fatal("error in setting the scocket to non-blocking mode: ", err)
	}

	// bind the socket to the ip and port
	ip := net.ParseIP(config.Config.Host)
	err = unix.Bind(serverFd, &unix.SockaddrInet4{
		Port: config.Config.Port,
		Addr: [4]byte(ip),
	})
	if err != nil {
		log.Fatal("error in binding the socket to the ip and port: ", err)
	}

	// start listening on the socket
	err = unix.Listen(serverFd, max_clients)
	if err != nil {
		log.Fatal("error in start listening on the socket", err)
	}

	log.Printf("Started listening on %s:%d", config.Config.Host, config.Config.Port)

	// create a kqueue
	kq, err := unix.Kqueue()
	defer unix.Close(kq)

	// create kqueue event for the server file descriptor
	serverKevent := unix.Kevent_t{
		Ident:  uint64(serverFd),
		Filter: unix.EVFILT_READ,
		Flags:  unix.EV_ADD | unix.EV_ENABLE,
	}

	// register the serverKevent with kqueue
	_, err = unix.Kevent(kq, []unix.Kevent_t{serverKevent}, nil, nil)
	if err != nil {
		log.Fatal("error in registering the server Kevent with kqueue ", err)
	}

	events := make([]unix.Kevent_t, max_clients)

	// wait for events
	for {

		n, err := unix.Kevent(kq, nil, events, nil)
		if err != nil {
			continue
		}

		for i := range n {
			if events[i].Ident == uint64(serverFd) {
				// new client wants to connect

				// accept the connection
				clientFd, sa, err := unix.Accept(serverFd)

				// create kevent for the new client
				clientKevent := unix.Kevent_t{
					Ident:  uint64(clientFd),
					Filter: unix.EVFILT_READ,
					Flags:  unix.EV_ADD | unix.EV_ENABLE,
				}

				// add the new client to the kqueue
				_, err = unix.Kevent(kq, []unix.Kevent_t{clientKevent}, nil, nil)
				if err != nil {
					log.Printf("Failed to accept the connection by: %v", sa)
					continue
				}

				con_clients++

			} else {
				// client wants to send data
				fdComm := core.FdComm{Fd: int(events[i].Ident)}
				cmd, err := readCommand(fdComm)
				if err != nil {
					log.Printf("disconnecting with: %v", err)
					unix.Close(fdComm.Fd)
					con_clients--
					continue
				}

				respond(fdComm, cmd)
			}
		}
	}
}
