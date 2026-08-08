package core

import (
	"errors"
	"strings"
)

// EvaluateResponse takes a redis command and
// evaluates the response and error(if any) for it
func EvaluateResponse(redisCmd RedisCmd) ([]byte, error) {

	cmd := strings.ToUpper(redisCmd.Command)
	args := redisCmd.Arguments

	switch cmd {
	case "PING":
		return evaluatePING(args)
	default:
		return evaluatePING(args)
	}
}

// evaluatePING evaluates the response for the
// redis "PING" command
func evaluatePING(args []string) ([]byte, error) {

	if len(args) > 1 {
		return nil, errors.New("ERR wrong number of arguments for the ping command")
	}

	var b []byte
	var err error

	if len(args) == 0 {
		b, err = Encode("PONG", true)
	} else {
		b, err = Encode(args[0], false)
	}

	if err != nil {
		return nil, err
	}

	return b, nil
}
