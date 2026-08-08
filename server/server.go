package server

import (
	"fmt"
	"log"
	"net"

	"github.com/cpd007/myredis/config"
	"github.com/cpd007/myredis/core"
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
			cmd, err := readCommand(conn)
			if err != nil {
				log.Printf("disconnecting with: %v", err)
				conn.Close()
				break
			}

			respond(conn, cmd)
		}
	}
}

func readCommand(c net.Conn) (cmd core.RedisCmd, err error) {

	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		return cmd, err
	}

	log.Printf("Read from connection: %s", string(buffer[:n]))

	decodedArray, err := core.DecodeStringArrays(buffer[:n])
	if err != nil {
		return cmd, err
	}

	return core.RedisCmd{
		Command:   decodedArray[0],
		Arguments: decodedArray[1:],
	}, nil
}

func respond(c net.Conn, cmd core.RedisCmd) {

	data, err := core.EvaluateResponse(cmd)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.Write(data)
}

func respondWithError(c net.Conn, err error) {

	c.Write([]byte(fmt.Sprintf("-%v\r\n", err)))
}
