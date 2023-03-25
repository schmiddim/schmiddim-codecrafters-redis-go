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
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage
	//
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
	// Make a buffer to hold incoming data.

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	brr := string(buf)
	fmt.Println("Message received:", brr)

	for _, element := range strings.Split(brr, "\n") {
		if strings.TrimRight(strings.ToLower(element), "\r") == "ping" {
			_, err = conn.Write([]byte("+PONG\r\n"))
		}
	}
	// Send a response back to person contacting us.
	if err != nil {
		return
	}
	// Close the connection when you're done with it.
	err = conn.Close()
	if err != nil {
		return
	}
}
