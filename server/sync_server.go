package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/cpd007/myredis/config"
)

func StartSyncTCPServer() {
	cfg := config.Config

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("error in intialising the listener", err)
	}

	defer listener.Close()

	log.Printf("TCP Server listening on %s\n", address)
	con_clients := 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error in accepting the connection", err)
		}

		con_clients++
		log.Printf("New connection established from %s\n", conn.RemoteAddr())

		for {
			cmd, err := readCommand(conn)
			if err != nil {
				log.Printf("disconnecting with: %v", err)
				conn.Close()
				con_clients--
				if err == io.EOF {
					break
				}
			}

			respond(conn, cmd)
		}
	}
}
