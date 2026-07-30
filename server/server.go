package server

import (
	"fmt"
	"log"
	"net"

	"github.com/cpd007/myredis/config"
)

// use command: <nc localhost 7379> to test.

func Start() {
	cfg := config.Config

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("error in intialising the listener", err)
	}

	defer listener.Close()

	log.Printf("TCP Server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error in accepting the connection", err)
		}

		log.Printf("New connection established from %s\n", conn.RemoteAddr())

		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				log.Printf("disconnecting with: %v", err)
				conn.Close()
				break
			}

			log.Printf("Read from connection: %s", string(buffer[:n]))

			// Write the same data back (Echo)
			conn.Write(buffer[:n])
		}
	}
}
