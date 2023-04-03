package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Type byte

const (
	simpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type CacheEntry struct {
	value       string
	expiryTime  int
	dateCreated time.Time
}

func NewCacheEntry(value string, expiryTime int) CacheEntry {
	return CacheEntry{
		dateCreated: time.Now(),
		value:       value,
		expiryTime:  expiryTime,
	}
}
func (c CacheEntry) String() string {
	return c.value
}
func (c CacheEntry) IsExpired() bool {
	if c.expiryTime < 0 {
		return false
	}
	duration := time.Duration(c.expiryTime) * time.Millisecond
	expiryDate := c.dateCreated.Add(duration)
	if expiryDate.Before(time.Now()) {
		return true
	}
	return false
}
func (c CacheEntry) Len() int {
	return len(c.value)
}

type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

func (v Value) String() string {
	if v.typ == BulkString || v.typ == simpleString {
		return strings.TrimSpace(string(v.bytes))
	}
	return ""

}

func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}
	return []Value{}

}

func decodeInput(byteStream *bufio.Reader) (Value, error) {

	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleStringString(byteStream)
	case "$":
		return decodeBulkString(byteStream)
	case "*":
		return decodeArray(byteStream)
	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleStringString(byteStream *bufio.Reader) (Value, error) {

	readBytes, err := readUntilCRLF(byteStream)

	if err != nil {
		return Value{}, err
	}

	return Value{typ: simpleString, bytes: readBytes}, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	var readBytes []byte
	for {

		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}
	return readBytes[:len(readBytes)-2], nil
}
func decodeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)

	if err != nil {
		return Value{}, err
	}
	count, err := strconv.Atoi(string(readBytesForCount))

	if err != nil {
		return Value{}, fmt.Errorf("unable to parse length: %s", err)
	}

	var array []Value
	for i := 1; i <= count; i++ {
		value, err := decodeInput(byteStream)
		if err != nil {
			return Value{}, err
		}
		array = append(array, value)
	}

	return Value{typ: Array, array: array}, nil
}
func decodeBulkString(byteStream *bufio.Reader) (Value, error) {

	readBytesForCount, err := readUntilCRLF(byteStream)

	if err != nil {
		return Value{}, err
	}
	count, err := strconv.Atoi(string(readBytesForCount))

	if err != nil {
		return Value{}, fmt.Errorf("unable to parse length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("unable to read: %s", err)

	}
	return Value{typ: BulkString, bytes: readBytes}, nil

}
