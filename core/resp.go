package core

import (
	"errors"
)

func readSimpleString(data []byte) (interface{}, int, error) {
	var pos int
	for pos = 1; pos < len(data); pos++ {
		if data[pos] == '\r' {
			break
		}
	}
	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (interface{}, int, error) {
	var pos int
	for pos = 1; pos < len(data); pos++ {
		if data[pos] == '\r' {
			break
		}
	}
	return string(data[1:pos]), pos + 2, nil
}

func readInt64(data []byte) (interface{}, int, error) {
	var pos int
	var value int64
	sign := int64(1)
	if len(data) > 1 && data[1] == '-' {
		sign = -1
		pos = 2
	} else {
		pos = 1
	}
	for ; pos < len(data); pos++ {
		if data[pos] == '\r' {
			break
		}
		value = value*10 + int64(data[pos]-'0')
	}
	return sign * value, pos + 2, nil
}

func readBulkString(data []byte) (interface{}, int, error) {
	var pos int
	var length int64
	sign := int64(1)
	if len(data) > 1 && data[1] == '-' {
		sign = -1
		pos = 2
	} else {
		pos = 1
	}
	for ; pos < len(data); pos++ {
		if data[pos] == '\r' {
			break
		}
		length = length*10 + int64(data[pos]-'0')
	}
	length *= sign
	if length == -1 {
		return nil, pos + 2, nil
	}
	pos += 2
	if pos+int(length)+2 > len(data) {
		return nil, 0, errors.New("incomplete bulk string")
	}
	return string(data[pos : pos+int(length)]),
		pos + int(length) + 2,
		nil
}

func readArray(data []byte) (interface{}, int, error) {
	var pos int
	var length int64
	for pos = 1; pos < len(data); pos++ {
		if data[pos] == '\r' {
			break
		}
		length = length*10 + int64(data[pos]-'0')
	}
	pos += 2
	elems := make([]interface{}, length)

	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	default:
		return nil, 0, errors.New("unknown RESP type")
	}
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}
