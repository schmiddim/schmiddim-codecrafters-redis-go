package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestResponse(t *testing.T) {
	tests := map[string]string{

		"pong": "+ping\r\n",
		//'$'
		//'*': "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n",
	}

	for assert, test := range tests {

		value, err := decodeInput(bufio.NewReader(bytes.NewBufferString(test)))
		if err != nil {
			t.Errorf("error decoding string: %s", err)

		}
		if (string(value.bytes)) != assert {
			t.Error("Bah", assert, value)

		}

	}
}
func TestType(t *testing.T) {

	tests := map[rune]string{

		'+': "+ping\r\n",
		'$': "$4\r\nECHO\r\n",
		'*': "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n",
	}

	for assert, test := range tests {

		value, err := decodeInput(bufio.NewReader(bytes.NewBufferString(test)))
		if err != nil {
			t.Errorf("error decoding string: %s", err)

		}
		if rune(value.typ) != assert {
			t.Error("Invalid Type ", value)

		}

	}

}
func TestFoo(t *testing.T) {
	//@todo implement me!
	//https://redis.io/docs/reference/protocol-spec/#send-commands-to-a-redis-server
	//https://redis.io/docs/reference/protocol-spec/
	//https://redis.io/commands/echo/
	//value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("+foo\r\n")))
	//if err != nil {
	//t.Errorf("error decoding simple string: %s", err)
	//}
	value, err := decodeInput(bufio.NewReader(bytes.NewBufferString("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")))

	if err != nil {
		t.Errorf("error decoding string: %s", err)
	}
	fmt.Println(value)
}
