package server

import (
	"fmt"
	"io"
	"log"

	"github.com/cpd007/myredis/core"
)

func readCommand(c io.ReadWriter) (cmd core.RedisCmd, err error) {

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

func respond(c io.ReadWriter, cmd core.RedisCmd) {

	data, err := core.EvaluateResponse(cmd)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.Write(data)
}

func respondWithError(c io.ReadWriter, err error) {

	c.Write([]byte(fmt.Sprintf("-%v\r\n", err)))
}
