package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	ConnHost = "0.0.0.0"
	ConnPort = "6379"
	ConnType = "tcp"
)

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
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		arr := strings.Split(string(buf), "\r\n")

		if string(buf[0]) == "*" {

			fmt.Println("is an array, size is", (string(buf[1])))

		}

		cmdFound := false
		for index, cmd := range arr {
			cmd = strings.TrimRight(strings.ToLower(cmd), "\r")
			if cmd == "ping" {
				_, err = conn.Write([]byte("+PONG\r\n"))
				cmdFound = true
				break
			}
			if cmd == "echo" {
				fmt.Println(index, arr[index+2])
				_, err = conn.Write([]byte("" + arr[index+2] + "\r\n"))

				cmdFound = true
				break
			}
		}

		if cmdFound == false {
			//	_, err = conn.Write([]byte("+command not found\r\n"))
			//
			fmt.Println("cmd not found")
			//return
		}

		//_, err = conn.Write([]byte("\r\n"))
		//if err != nil {
		//	return
		//}
	}
}
