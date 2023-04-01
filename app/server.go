package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	ConnHost = "0.0.0.0"
	ConnPort = "6379"
	ConnType = "tcp"
)

var cacheItems = make(map[string]string)

func main() {
	l, err := net.Listen(ConnType, ConnHost+":"+ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {

		}
	}(l)
	fmt.Println("Listening on " + ConnHost + ":" + ConnPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		go handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection")
		}
	}(conn)
	for {
		value, err := decodeInput(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Print("Err reading client ", err.Error())
			return
		}
		command := value.Array()[0].String()
		args := value.Array()[1:]

		switch command {
		case "ping":
			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				return
			}

		case "echo":
			_, err := conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
			if err != nil {
				return
			}
		case "set":
			if len(value.Array()) != 3 {
				_, err := conn.Write([]byte("-ERR your doing it wrong\r\n"))
				if err != nil {
					return
				}

			}
			key := value.Array()[1].String()
			cacheItems[key] = value.Array()[2].String()
			_, err := conn.Write([]byte("+saved\r\n"))
			if err != nil {
				return
			}
		case "get":
			if len(value.Array()) != 2 {
				_, err := conn.Write([]byte("-ERR your doing it wrong\r\n"))
				if err != nil {
					return
				}
			}
			key := value.Array()[1].String()
			n := strconv.Itoa(len(cacheItems[key]))
			item := cacheItems[key]

			resultString := "$" + n + "\r\n" + item + "\r\n"
			stream := []byte(resultString)

			_, err := conn.Write(stream)
			if err != nil {
				return
			}
		default:
			_, err2 := conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
			if err2 != nil {
				return
			}
		}

	}

}
