package core

import (
	"errors"
	"fmt"
)

// RESP decoder functions

// DecodeStringArrays decodes the byte stream for resp protocol
// into a string array
func DecodeStringArrays(data []byte) ([]string, error) {

	d, err := Decode(data)
	if err != nil {
		return nil, err
	}

	ds, ok := d.([]any)
	if !ok {
		return nil, errors.New("Unsupported arguments entered")
	}

	stringArr := make([]string, 0)

	for _, v := range ds {
		str, ok := v.(string)
		if !ok {
			return nil, errors.New("Unsupported arguments entered")
		}
		stringArr = append(stringArr, str)
	}

	return stringArr, nil
}

// Decode decodes the byte stream for resp protocol
func Decode(data []byte) (any, error) {

	if len(data) == 0 {
		return nil, errors.New("Empty data received.")
	}

	value, _, err := decodeOne(data)
	return value, err
}

// decodeOne decodes data one at a time and returns
// the decoded data, the delta it has decoded till now, and error
func decodeOne(data []byte) (any, int, error) {

	indicator := string(data[0])

	switch indicator {
	case "+":
		return readSimpleString(data)
	case "-":
		return readError(data)
	case ":":
		return readInteger(data)
	case "$":
		return readBulkStrings(data)
	case "*":
		return readArray(data)
	default:
		return nil, 0, errors.New("Unsupported command entered.")
	}
}

// readSimpleString parses a RESP simple string reply.
// It returns the string value, the delta it has decoded till now, and error.
func readSimpleString(data []byte) (string, int, error) {

	// first character is "+"
	pos := 1

	for ; pos < len(data) && data[pos] != '\r'; pos++ {
	}

	if pos == len(data) {
		return "", 0, errors.New("Invalid format for simple string")
	}

	return string(data[1:pos]), pos + 2, nil
}

// readError parses a RESP error reply.
// It returns the error text, the delta it has decoded till now, and error.
func readError(data []byte) (string, int, error) {

	val, delta, err := readSimpleString(data)
	if err != nil {
		return "", 0, errors.New("Invalid format for error")
	}

	return val, delta, nil
}

// readInteger parses a RESP integer reply.
// It returns the value, the delta it has decoded till now, and error.
func readInteger(data []byte) (int64, int, error) {

	val, delta, err := readNumber(data[1:])
	if err != nil {
		return 0, 0, err
	}

	return val, delta + 1, nil
}

// readBulkStrings parses a RESP bulk string reply.
// It returns the bulk string value, the delta it has decoded till now, and error.
func readBulkStrings(data []byte) (string, int, error) {
	if len(data) == 0 || data[0] != '$' {
		return "", 0, errors.New("Invalid format for bulk string length")
	}

	length, delta, err := readNumber(data[1:])
	if err != nil {
		return "", 0, errors.New("Invalid format for bulk string length")
	}

	pos := delta + 1
	total := pos + int(length)

	if total+2 > len(data) {
		return "", 0, errors.New("Invalid format for bulk string payload")
	}

	return string(data[pos:total]), total + 2, nil
}

// readArray parses a RESP array reply.
// It returns the elements, the delta it has decoded till now, and error.
func readArray(data []byte) ([]any, int, error) {

	length, delta, err := readNumber(data[1:])
	if err != nil {
		return nil, 0, errors.New("Invalid format for array length")
	}

	pos := delta + 1
	elems := make([]any, length)

	for i := range length {
		elems[i], delta, err = decodeOne(data[pos:])
		if err != nil {
			return nil, 0, errors.Join(err, errors.New("Invalid format for array"))
		}

		pos += delta
	}

	return elems, pos, nil
}

// readNumber parses a RESP integer reply without the indicator.
// It returns the number, the delta it has decoded till now, and error.
func readNumber(data []byte) (int64, int, error) {
	pos := 0
	val := int64(0)

	for ; pos < len(data) && data[pos] != '\r'; pos++ {
		if data[pos] < '0' || data[pos] > '9' {
			return 0, 0, errors.New("Invalid format for integer")
		}
		val = 10*val + int64(data[pos]-'0')
	}

	if pos == len(data) {
		return 0, 0, errors.New("Invalid format for integer")
	}

	return val, pos + 2, nil
}

// RESP encoder functions

// Encode function encodes the given response value
// to RESP response according to its type
func Encode(value any, isSimple bool) ([]byte, error) {

	switch v := value.(type) {
	case string:
		if isSimple {
			return encodeSimpleString(v)
		}
		return encodeBulkString(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return encodeInteger(v)
	}
	return []byte{}, nil
}

// encodeSimpleString encodes a string
// into RESP's simple string
func encodeSimpleString(s string) ([]byte, error) {

	encStr := fmt.Sprintf("+%s\r\n", s)
	return []byte(encStr), nil
}

// encodeBulkString encodes a string
// into RESP's bulk string
func encodeBulkString(s string) ([]byte, error) {

	encStr := fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
	return []byte(encStr), nil
}

// encodeInteger encodes an integer
// into RESP's integer
func encodeInteger(v any) ([]byte, error) {

	encInt := fmt.Sprintf(":%d\r\n", v)
	return []byte(encInt), nil
}

// encodeNil encodes a RESP nil response
func encodeNil() []byte {

	return []byte("$-1\r\n")
}
