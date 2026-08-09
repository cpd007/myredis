package core

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// EvaluateResponse takes a redis command and
// evaluates the response and error(if any) for it
func EvaluateResponse(redisCmd RedisCmd) ([]byte, error) {

	cmd := strings.ToUpper(redisCmd.Command)
	args := redisCmd.Arguments

	switch cmd {
	case "PING":
		return evaluatePING(args)
	case "SET":
		return evaluateSET(args)
	case "GET":
		return evaluateGET(args)
	case "TTL":
		return evaluateTTL(args)
	default:
		return evaluatePING(args)
	}
}

// evaluatePING evaluates the response for the
// redis "PING" command
func evaluatePING(args []string) ([]byte, error) {

	if len(args) > 1 {
		return nil, errors.New("ERR wrong number of arguments for the 'ping' command")
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

// evaluateSET sets the value of the specified
// key in the hash table with a TTL(if present)
func evaluateSET(args []string) ([]byte, error) {

	if len(args) < 2 {
		return nil, errors.New("ERR wrong number of arguments for the 'set' command")
	}

	exDurationMs := int64(-1)
	key := args[0]
	value := args[1]

	for i := 2; i < len(args); i++ {

		switch strings.ToUpper(args[i]) {
		case "TTL":
			i++
			if i == len(args) {
				return nil, errors.New("ERR wrong number of arguments for the 'set' command")
			}

			exDurationSec, err := strconv.ParseUint(args[i], 10, 32)
			if err != nil {
				return nil, errors.New("ERR invalid redis syntax")
			}

			exDurationMs = int64(exDurationSec * 1000)
		default:
			return nil, errors.New("ERR invalid redis syntax")
		}
	}

	Put(key, NewObj(value, exDurationMs))

	return Encode("OK", false)
}

// evaluateGET gets the value (if present) for the specified kwy
func evaluateGET(args []string) ([]byte, error) {

	if len(args) > 1 {
		return nil, errors.New("ERR wrong number of arguments for the 'get' command")
	}

	obj := Get(args[0])

	// no key present
	if obj == nil {
		return encodeNil(), nil
	}
	// key expired, no key present
	if obj.ExpiresAt > 0 && time.Now().UnixMilli() > obj.ExpiresAt {
		return encodeNil(), nil
	}

	return Encode(obj.Value, false)
}

// evaluateTTL gets the TTL (if present) for the specified kwy
func evaluateTTL(args []string) ([]byte, error) {

	if len(args) > 1 {
		return nil, errors.New("ERR wrong number of arguments for the 'ttl' command")
	}

	obj := Get(args[0])

	// no key present
	if obj == nil {
		return Encode(-2, false)
	}
	// no expiry
	if obj.ExpiresAt == -1 {
		return Encode(-1, false)
	}
	// key expired, no key present
	if obj.ExpiresAt > 0 && time.Now().UnixMilli() > obj.ExpiresAt {
		return Encode(-2, false)
	}

	return Encode((obj.ExpiresAt-time.Now().UnixMilli())/1000, false)
}
