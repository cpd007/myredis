package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cpd007/myredis/core"
	"github.com/cpd007/myredis/helper"
)

func readCommands(c io.ReadWriter) (cmd []core.RedisCmd, err error) {

	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		return cmd, err
	}

	log.Printf("Read from connection: %s", string(buffer[:n]))

	decodedValues, err := core.Decode(buffer[:n])
	if err != nil {
		return cmd, err
	}

	cmds := []core.RedisCmd{}

	for i := range decodedValues {

		decodedStringArray, err := helper.ToStringArray(decodedValues[i])
		if err != nil {
			log.Println(err)
			return nil, errors.New("Invalid request")
		}

		cmds = append(cmds, core.RedisCmd{
			Command:   strings.ToUpper(decodedStringArray[0]),
			Arguments: decodedStringArray[1:],
		})
	}

	return cmds, nil
}

func respond(c io.ReadWriter, cmds []core.RedisCmd) {

	var response []byte

	for i := range cmds {

		data, err := core.EvaluateResponse(cmds[i])
		if err != nil {
			data = respondWithError(err)
		}

		response = append(response, data...)
	}

	c.Write(response)
}

func respondWithError(err error) []byte {

	return fmt.Appendf(nil, "-%v\r\n", err)
}
